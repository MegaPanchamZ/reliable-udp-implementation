package protocol

import (
	"urp-go/internal/logger"
	"urp-go/internal/plc"
	"urp-go/internal/urp"
	"fmt"
	"net"
	"os"
	"time"
)

// Receiver implements the receiver-side URP protocol with 4-state machine
type Receiver struct {
	conn       *net.UDPConn
	senderAddr *net.UDPAddr
	logger     *logger.Logger
	plc        *plc.Module
	state      int
	maxWindow  int

	// Receive buffer for out-of-order segments
	buffer      map[uint16]*urp.URPSegment
	expectedSeq uint16

	// Output file
	outputFile *os.File

	// Channels
	done chan bool
}

// NewReceiver creates a new Receiver instance
func NewReceiver(port int, senderPort int, outputPath string, maxWindow int, log *logger.Logger, plcModule *plc.Module) (*Receiver, error) {
	// Create UDP address
	addr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
	}

	// Create UDP connection
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP connection: %w", err)
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	r := &Receiver{
		conn:        conn,
		logger:      log,
		plc:         plcModule,
		state:       urp.StateLISTEN,
		maxWindow:   maxWindow,
		buffer:      make(map[uint16]*urp.URPSegment),
		expectedSeq: 0,
		outputFile:  outputFile,
		done:        make(chan bool),
	}

	return r, nil
}

// Close closes the receiver connection and output file
func (r *Receiver) Close() {
	close(r.done)
	r.conn.Close()
	if r.outputFile != nil {
		r.outputFile.Close()
	}
}

// Listen starts listening for incoming segments
func (r *Receiver) Listen() error {
	buffer := make([]byte, urp.HeaderSize+urp.MSS)

	for {
		select {
		case <-r.done:
			return nil
		default:
			// Read incoming segment
			n, addr, err := r.conn.ReadFromUDP(buffer)
			if err != nil {
				continue
			}

			// Remember sender address
			if r.senderAddr == nil {
				r.senderAddr = addr
			}

			// Unpack segment
			seg, err := urp.Unpack(buffer[:n])
			if err != nil {
				continue
			}

			// Validate checksum
			if !seg.IsValid() {
				// Drop corrupted packet
				continue
			}

			// Log received segment
			r.logger.LogSegment(logger.EventReceive, seg.SeqNum, seg.AckNum,
				r.flagsToString(seg.Flags), len(seg.Payload))

			// Process segment based on state
			if err := r.processSegment(seg); err != nil {
				return err
			}
		}
	}
}

// processSegment processes an incoming segment based on current state
func (r *Receiver) processSegment(seg *urp.URPSegment) error {
	switch r.state {
	case urp.StateLISTEN:
		return r.handleListen(seg)
	case urp.StateSYNRCVD:
		return r.handleSynReceived(seg)
	case urp.StateESTABLISHED:
		return r.handleEstablished(seg)
	case urp.StateTIMEWAIT:
		return r.handleTimeWait(seg)
	}
	return nil
}

// handleListen handles segments in LISTEN state
func (r *Receiver) handleListen(seg *urp.URPSegment) error {
	if seg.HasFlag(urp.FlagSYN) {
		r.changeState(urp.StateSYNRCVD)
		r.expectedSeq = seg.SeqNum + 1

		// Send SYN-ACK
		ackSeg := urp.NewSegment(0, r.expectedSeq, urp.FlagSYN|urp.FlagACK, nil)
		return r.sendACK(ackSeg)
	}
	return nil
}

// handleSynReceived handles segments in SYN_RCVD state
func (r *Receiver) handleSynReceived(seg *urp.URPSegment) error {
	if seg.HasFlag(urp.FlagACK) {
		r.changeState(urp.StateESTABLISHED)
		return nil
	}

	// If we receive data, transition to ESTABLISHED
	if len(seg.Payload) > 0 {
		r.changeState(urp.StateESTABLISHED)
		return r.handleEstablished(seg)
	}

	return nil
}

// handleEstablished handles segments in ESTABLISHED state
func (r *Receiver) handleEstablished(seg *urp.URPSegment) error {
	// Handle FIN
	if seg.HasFlag(urp.FlagFIN) {
		r.changeState(urp.StateTIMEWAIT)

		// Send FIN-ACK
		ackSeg := urp.NewSegment(0, seg.SeqNum+1, urp.FlagACK, nil)
		r.sendACK(ackSeg)

		// Start TIME_WAIT timer
		go func() {
			time.Sleep(time.Duration(urp.TimeWaitDuration) * time.Millisecond)
			r.done <- true
		}()

		return nil
	}

	// Handle data segment
	if len(seg.Payload) > 0 {
		return r.handleDataSegment(seg)
	}

	return nil
}

// handleTimeWait handles segments in TIME_WAIT state
func (r *Receiver) handleTimeWait(seg *urp.URPSegment) error {
	// In TIME_WAIT, just acknowledge any retransmitted FINs
	if seg.HasFlag(urp.FlagFIN) {
		ackSeg := urp.NewSegment(0, seg.SeqNum+1, urp.FlagACK, nil)
		return r.sendACK(ackSeg)
	}
	return nil
}

// handleDataSegment processes a data segment
func (r *Receiver) handleDataSegment(seg *urp.URPSegment) error {
	// Check if this is the expected segment
	if seg.SeqNum == r.expectedSeq {
		// Write to file
		if _, err := r.outputFile.Write(seg.Payload); err != nil {
			return fmt.Errorf("failed to write to output file: %w", err)
		}

		r.expectedSeq += uint16(len(seg.Payload))

		// Check if any buffered segments can now be delivered
		for {
			if bufferedSeg, exists := r.buffer[r.expectedSeq]; exists {
				r.outputFile.Write(bufferedSeg.Payload)
				r.expectedSeq += uint16(len(bufferedSeg.Payload))
				delete(r.buffer, bufferedSeg.SeqNum)
			} else {
				break
			}
		}
	} else if seg.SeqNum > r.expectedSeq {
		// Out-of-order segment - buffer it
		r.buffer[seg.SeqNum] = seg
	}
	// If seg.SeqNum < r.expectedSeq, it's a duplicate - just acknowledge

	// Always send ACK with the next expected sequence number
	ackSeg := urp.NewSegment(0, r.expectedSeq, urp.FlagACK, nil)
	return r.sendACK(ackSeg)
}

// sendACK sends an ACK segment to the sender
func (r *Receiver) sendACK(seg *urp.URPSegment) error {
	if r.senderAddr == nil {
		return fmt.Errorf("sender address not known")
	}

	data := seg.Pack()

	// Apply PLC for reverse path (ACKs)
	if r.plc != nil {
		processedData, shouldSend := r.plc.ProcessOutgoingReverse(data, seg)
		if !shouldSend {
			// Packet dropped by PLC
			return nil
		}
		data = processedData
	}

	_, err := r.conn.WriteToUDP(data, r.senderAddr)
	if err != nil {
		return err
	}

	r.logger.LogSegment(logger.EventSend, seg.SeqNum, seg.AckNum,
		r.flagsToString(seg.Flags), len(seg.Payload))

	return nil
}

// changeState changes the receiver state and logs it
func (r *Receiver) changeState(newState int) {
	oldState := r.state
	r.state = newState
	r.logger.LogStateChange(r.stateToString(oldState), r.stateToString(newState))
}

// stateToString converts state to string
func (r *Receiver) stateToString(state int) string {
	switch state {
	case urp.StateLISTEN:
		return "LISTEN"
	case urp.StateSYNRCVD:
		return "SYN_RCVD"
	case urp.StateESTABLISHED:
		return "ESTABLISHED"
	case urp.StateTIMEWAIT:
		return "TIME_WAIT"
	default:
		return "UNKNOWN"
	}
}

// flagsToString converts flags to string representation
func (r *Receiver) flagsToString(flags uint8) string {
	result := ""
	if flags&urp.FlagACK != 0 {
		result += "ACK"
	}
	if flags&urp.FlagSYN != 0 {
		if result != "" {
			result += "|"
		}
		result += "SYN"
	}
	if flags&urp.FlagFIN != 0 {
		if result != "" {
			result += "|"
		}
		result += "FIN"
	}
	if result == "" {
		result = "NONE"
	}
	return result
}

package protocol

import (
	"fmt"
	"net"
	"time"
	"urp-go/internal/logger"
	"urp-go/internal/plc"
	"urp-go/internal/urp"
)

// Sender implements the sender-side URP protocol with 5-state machine
type Sender struct {
	conn       *net.UDPConn
	remoteAddr *net.UDPAddr
	logger     *logger.Logger
	plc        *plc.Module
	state      int
	maxWindow  int
	rto        time.Duration

	// Sliding window
	sendBase   uint16
	nextSeqNum uint16
	window     map[uint16]*urp.URPSegment

	// Timers
	timer       *time.Timer
	timerActive bool

	// Fast retransmit
	dupACKCount map[uint16]int

	// Channels
	ackChan chan *urp.URPSegment
	done    chan bool
}

// NewSender creates a new Sender instance
func NewSender(localPort int, remoteHost string, remotePort int, maxWindow int, rto int,
	log *logger.Logger, plcModule *plc.Module) (*Sender, error) {

	// Resolve local address
	localAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: localPort,
	}

	// Resolve remote address
	remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", remoteHost, remotePort))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve remote address: %w", err)
	}

	// Create UDP connection
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP connection: %w", err)
	}

	s := &Sender{
		conn:        conn,
		remoteAddr:  remoteAddr,
		logger:      log,
		plc:         plcModule,
		state:       urp.StateCLOSED,
		maxWindow:   maxWindow,
		rto:         time.Duration(rto) * time.Millisecond,
		sendBase:    0,
		nextSeqNum:  0,
		window:      make(map[uint16]*urp.URPSegment),
		dupACKCount: make(map[uint16]int),
		ackChan:     make(chan *urp.URPSegment, 100),
		done:        make(chan bool),
	}

	// Start ACK receiver goroutine
	go s.receiveACKs()

	return s, nil
}

// Close closes the sender connection
func (s *Sender) Close() {
	close(s.done)
	s.conn.Close()
}

// receiveACKs listens for incoming ACK segments
func (s *Sender) receiveACKs() {
	buffer := make([]byte, urp.HeaderSize+urp.MSS)

	for {
		select {
		case <-s.done:
			return
		default:
			s.conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
			n, _, err := s.conn.ReadFromUDP(buffer)
			if err != nil {
				continue // Timeout or other error
			}

			// Process through PLC (reverse direction)
			processedData, status, shouldDeliver := s.plc.ProcessIncomingReverse(buffer[:n])

			// Log the received segment (even if dropped or corrupted)
			seg, err := urp.Unpack(buffer[:n])
			if err == nil {
				segType := urp.GetSegmentType(seg)
				s.logger.LogSegment(logger.EventReceive, status, segType, seg.SeqNum, len(seg.Payload))
			}

			if !shouldDeliver {
				continue // Dropped by PLC
			}

			// Unpack the processed data
			seg, err = urp.Unpack(processedData)
			if err != nil {
				continue
			}

			// Validate checksum
			if !seg.IsValid() {
				s.logger.IncrementStat("corrupt_acks_discarded", 1)
				continue // Drop corrupted packet
			}

			s.ackChan <- seg
		}
	}
}

// SendFile sends a file using the URP protocol
func (s *Sender) SendFile(data []byte) error {
	// State: CLOSED -> SYN_SENT
	if err := s.sendSYN(); err != nil {
		return err
	}

	// Wait for SYN-ACK
	if err := s.waitForSYNACK(); err != nil {
		return err
	}

	// State: ESTABLISHED
	// Send data segments
	if err := s.sendData(data); err != nil {
		return err
	}

	// State: ESTABLISHED -> FIN_SENT
	if err := s.sendFIN(); err != nil {
		return err
	}

	// Wait for FIN-ACK
	if err := s.waitForFINACK(); err != nil {
		return err
	}

	// Write statistics to log file
	s.logger.WriteSenderStatistics()

	return nil
}

// sendSYN sends the initial SYN segment
func (s *Sender) sendSYN() error {
	s.changeState(urp.StateSYNSENT)

	seg := urp.NewSegment(s.nextSeqNum, urp.FlagSYN, nil)
	s.logger.IncrementStat("original_segments_sent", 1)
	return s.sendSegment(seg, true)
}

// waitForSYNACK waits for the SYN-ACK response
func (s *Sender) waitForSYNACK() error {
	timeout := time.After(s.rto * 10) // Give more time for initial handshake

	for {
		select {
		case seg := <-s.ackChan:
			// ACK segments use SeqNum field as acknowledgment number
			if seg.HasFlag(urp.FlagACK) {
				s.sendBase = seg.SeqNum
				s.nextSeqNum = seg.SeqNum
				s.changeState(urp.StateESTABLISHED)
				return nil
			}
		case <-timeout:
			return fmt.Errorf("timeout waiting for SYN-ACK")
		}
	}
}

// sendData sends file data using sliding window
func (s *Sender) sendData(data []byte) error {
	// Split data into segments
	offset := 0
	for offset < len(data) {
		// Wait if window is full
		for s.getWindowSize() >= s.maxWindow {
			select {
			case seg := <-s.ackChan:
				s.handleACK(seg)
			case <-time.After(s.rto):
				s.handleTimeout()
			}
		}

		// Send next segment
		end := offset + urp.MSS
		if end > len(data) {
			end = len(data)
		}

		payload := data[offset:end]
		seg := urp.NewSegment(s.nextSeqNum, 0, payload) // 0 flags = DATA segment

		s.window[s.nextSeqNum] = seg

		// Track original data statistics
		s.logger.IncrementStat("original_data_sent", len(payload))
		s.logger.IncrementStat("original_segments_sent", 1)

		s.sendSegment(seg, true)

		s.nextSeqNum += uint16(len(payload))
		offset = end

		// Start timer if not active
		if !s.timerActive {
			s.startTimer()
		}
	}

	// Wait for all ACKs
	for len(s.window) > 0 {
		select {
		case seg := <-s.ackChan:
			s.handleACK(seg)
		case <-time.After(s.rto):
			s.handleTimeout()
		}
	}

	return nil
}

// sendFIN sends the FIN segment
func (s *Sender) sendFIN() error {
	s.changeState(urp.StateFINWAIT)

	seg := urp.NewSegment(s.nextSeqNum, urp.FlagFIN, nil)
	s.logger.IncrementStat("original_segments_sent", 1)
	return s.sendSegment(seg, true)
}

// waitForFINACK waits for the FIN-ACK response
func (s *Sender) waitForFINACK() error {
	timeout := time.After(s.rto * 10)

	for {
		select {
		case seg := <-s.ackChan:
			if seg.HasFlag(urp.FlagACK) {
				s.changeState(urp.StateCLOSED)
				return nil
			}
		case <-timeout:
			return fmt.Errorf("timeout waiting for FIN-ACK")
		}
	}
}

// sendSegment sends a segment through the PLC module
func (s *Sender) sendSegment(seg *urp.URPSegment, trackStats bool) error {
	data := seg.Pack()

	// Process through PLC
	processedData, status, shouldSend := s.plc.ProcessOutgoingForward(data, seg)

	// Always log the segment (even if dropped or corrupted)
	segType := urp.GetSegmentType(seg)
	s.logger.LogSegment(logger.EventSend, status, segType, seg.SeqNum, len(seg.Payload))

	// Track statistics
	if trackStats {
		if segType == logger.SegmentDATA {
			s.logger.IncrementStat("total_data_sent", len(seg.Payload))
		}
		s.logger.IncrementStat("total_segments_sent", 1)
	}

	if shouldSend && processedData != nil {
		_, err := s.conn.WriteToUDP(processedData, s.remoteAddr)
		if err != nil {
			return err
		}
	}

	return nil
}

// handleACK processes an incoming ACK
func (s *Sender) handleACK(seg *urp.URPSegment) {
	if !seg.HasFlag(urp.FlagACK) {
		return
	}

	// For ACK segments, SeqNum field contains the acknowledgment number
	ackNum := seg.SeqNum

	// Check for duplicate ACK
	if ackNum == s.sendBase {
		s.dupACKCount[ackNum]++
		s.logger.IncrementStat("duplicate_acks_received", 1)

		// Fast retransmit on 3 duplicate ACKs
		if s.dupACKCount[ackNum] >= 3 {
			if seg, exists := s.window[s.sendBase]; exists {
				s.logger.IncrementStat("fast_retransmits", 1)
				s.sendSegment(seg, false)
			}
			s.dupACKCount[ackNum] = 0
		}
		return
	}

	// Cumulative ACK - remove acknowledged segments
	for seqNum := s.sendBase; seqNum < ackNum; seqNum++ {
		delete(s.window, seqNum)
	}

	s.sendBase = ackNum

	// Reset duplicate ACK count
	s.dupACKCount = make(map[uint16]int)

	// Restart timer if window not empty
	if len(s.window) > 0 {
		s.startTimer()
	} else {
		s.stopTimer()
	}
}

// handleTimeout handles retransmission timeout
func (s *Sender) handleTimeout() {
	s.logger.IncrementStat("timeout_retransmits", 1)

	// Retransmit the oldest unacknowledged segment
	if seg, exists := s.window[s.sendBase]; exists {
		s.sendSegment(seg, false)
	}

	// Restart timer
	s.startTimer()
}

// getWindowSize returns the current window size in bytes
func (s *Sender) getWindowSize() int {
	size := 0
	for _, seg := range s.window {
		size += len(seg.Payload)
	}
	return size
}

// startTimer starts or restarts the retransmission timer
func (s *Sender) startTimer() {
	if s.timer != nil {
		s.timer.Stop()
	}
	s.timer = time.NewTimer(s.rto)
	s.timerActive = true

	go func() {
		<-s.timer.C
		if s.timerActive {
			s.handleTimeout()
		}
	}()
}

// stopTimer stops the retransmission timer
func (s *Sender) stopTimer() {
	if s.timer != nil {
		s.timer.Stop()
	}
	s.timerActive = false
}

// changeState changes the sender state
func (s *Sender) changeState(newState int) {
	s.state = newState
}

// stateToString converts state to string
func (s *Sender) stateToString(state int) string {
	switch state {
	case urp.StateCLOSED:
		return "CLOSED"
	case urp.StateSYNSENT:
		return "SYN_SENT"
	case urp.StateESTABLISHED:
		return "ESTABLISHED"
	case urp.StateCLOSING:
		return "CLOSING"
	case urp.StateFINWAIT:
		return "FIN_WAIT"
	default:
		return "UNKNOWN"
	}
}

// flagsToString converts flags to string representation
func (s *Sender) flagsToString(flags uint8) string {
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

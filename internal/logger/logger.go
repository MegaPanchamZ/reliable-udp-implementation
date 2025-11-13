package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// EventType represents the type of event being logged
type EventType string

const (
	EventSend    EventType = "snd"
	EventReceive EventType = "rcv"
)

// StatusType represents the PLC status of a segment
type StatusType string

const (
	StatusOK   StatusType = "ok"
	StatusDrop StatusType = "drp"
	StatusCor  StatusType = "cor"
)

// SegmentType represents the type of URP segment
type SegmentType string

const (
	SegmentDATA SegmentType = "DATA"
	SegmentACK  SegmentType = "ACK"
	SegmentSYN  SegmentType = "SYN"
	SegmentFIN  SegmentType = "FIN"
)

// Statistics tracks protocol statistics
type Statistics struct {
	// URP Protocol statistics
	OriginalDataSent      int
	TotalDataSent         int
	OriginalSegmentsSent  int
	TotalSegmentsSent     int
	TimeoutRetransmits    int
	FastRetransmits       int
	DuplicateACKsReceived int
	CorruptACKsDiscarded  int

	// PLC statistics
	PLCForwardDropped   int
	PLCForwardCorrupted int
	PLCReverseDropped   int
	PLCReverseCorrupted int

	// Receiver statistics
	OriginalDataReceived      int
	TotalDataReceived         int
	OriginalSegmentsReceived  int
	TotalSegmentsReceived     int
	CorruptSegmentsDiscarded  int
	DuplicateSegmentsReceived int
	TotalACKsSent             int
	DuplicateACKsSent         int
}

// Logger provides thread-safe logging to a file per spec format
type Logger struct {
	file  *os.File
	mu    sync.Mutex
	start time.Time
	stats Statistics
}

// New creates a new Logger that writes to the specified file
func New(filename string) (*Logger, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	return &Logger{
		file:  file,
		start: time.Now(),
	}, nil
}

// Close closes the log file
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}

// getTimestamp returns the elapsed time in milliseconds since logger creation
func (l *Logger) getTimestamp() float64 {
	elapsed := time.Since(l.start)
	return float64(elapsed.Microseconds()) / 1000.0
}

// LogSegment logs a segment per spec format:
// <direction> <status> <time> <segment-type> <seq-number> <payload-length>
func (l *Logger) LogSegment(direction EventType, status StatusType, segType SegmentType, seqNum uint16, payloadLen int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := l.getTimestamp()
	line := fmt.Sprintf("%s %s %.2f %s %d %d\n",
		direction, status, timestamp, segType, seqNum, payloadLen)
	l.file.WriteString(line)
}

// IncrementStat increments a statistic counter
func (l *Logger) IncrementStat(stat string, value int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	switch stat {
	case "original_data_sent":
		l.stats.OriginalDataSent += value
	case "total_data_sent":
		l.stats.TotalDataSent += value
	case "original_segments_sent":
		l.stats.OriginalSegmentsSent += value
	case "total_segments_sent":
		l.stats.TotalSegmentsSent += value
	case "timeout_retransmits":
		l.stats.TimeoutRetransmits += value
	case "fast_retransmits":
		l.stats.FastRetransmits += value
	case "duplicate_acks_received":
		l.stats.DuplicateACKsReceived += value
	case "corrupt_acks_discarded":
		l.stats.CorruptACKsDiscarded += value
	case "plc_forward_dropped":
		l.stats.PLCForwardDropped += value
	case "plc_forward_corrupted":
		l.stats.PLCForwardCorrupted += value
	case "plc_reverse_dropped":
		l.stats.PLCReverseDropped += value
	case "plc_reverse_corrupted":
		l.stats.PLCReverseCorrupted += value
	case "original_data_received":
		l.stats.OriginalDataReceived += value
	case "total_data_received":
		l.stats.TotalDataReceived += value
	case "original_segments_received":
		l.stats.OriginalSegmentsReceived += value
	case "total_segments_received":
		l.stats.TotalSegmentsReceived += value
	case "corrupt_segments_discarded":
		l.stats.CorruptSegmentsDiscarded += value
	case "duplicate_segments_received":
		l.stats.DuplicateSegmentsReceived += value
	case "total_acks_sent":
		l.stats.TotalACKsSent += value
	case "duplicate_acks_sent":
		l.stats.DuplicateACKsSent += value
	}
}

// WriteSenderStatistics writes sender statistics to log file per spec
func (l *Logger) WriteSenderStatistics() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.file.WriteString("\n")
	l.file.WriteString(fmt.Sprintf("Original data sent: %d\n", l.stats.OriginalDataSent))
	l.file.WriteString(fmt.Sprintf("Total data sent: %d\n", l.stats.TotalDataSent))
	l.file.WriteString(fmt.Sprintf("Original segments sent: %d\n", l.stats.OriginalSegmentsSent))
	l.file.WriteString(fmt.Sprintf("Total segments sent: %d\n", l.stats.TotalSegmentsSent))
	l.file.WriteString(fmt.Sprintf("Timeout retransmissions: %d\n", l.stats.TimeoutRetransmits))
	l.file.WriteString(fmt.Sprintf("Fast retransmissions: %d\n", l.stats.FastRetransmits))
	l.file.WriteString(fmt.Sprintf("Duplicate acks received: %d\n", l.stats.DuplicateACKsReceived))
	l.file.WriteString(fmt.Sprintf("Corrupted acks discarded: %d\n", l.stats.CorruptACKsDiscarded))
	l.file.WriteString(fmt.Sprintf("PLC forward segments dropped: %d\n", l.stats.PLCForwardDropped))
	l.file.WriteString(fmt.Sprintf("PLC forward segments corrupted: %d\n", l.stats.PLCForwardCorrupted))
	l.file.WriteString(fmt.Sprintf("PLC reverse segments dropped: %d\n", l.stats.PLCReverseDropped))
	l.file.WriteString(fmt.Sprintf("PLC reverse segments corrupted: %d\n", l.stats.PLCReverseCorrupted))
}

// WriteReceiverStatistics writes receiver statistics to log file per spec
func (l *Logger) WriteReceiverStatistics() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.file.WriteString("\n")
	l.file.WriteString(fmt.Sprintf("Original data received: %d\n", l.stats.OriginalDataReceived))
	l.file.WriteString(fmt.Sprintf("Total data received: %d\n", l.stats.TotalDataReceived))
	l.file.WriteString(fmt.Sprintf("Original segments received: %d\n", l.stats.OriginalSegmentsReceived))
	l.file.WriteString(fmt.Sprintf("Total segments received: %d\n", l.stats.TotalSegmentsReceived))
	l.file.WriteString(fmt.Sprintf("Corrupted segments discarded: %d\n", l.stats.CorruptSegmentsDiscarded))
	l.file.WriteString(fmt.Sprintf("Duplicate segments received: %d\n", l.stats.DuplicateSegmentsReceived))
	l.file.WriteString(fmt.Sprintf("Total acks sent: %d\n", l.stats.TotalACKsSent))
	l.file.WriteString(fmt.Sprintf("Duplicate acks sent: %d\n", l.stats.DuplicateACKsSent))
}

// Flush ensures all buffered data is written to disk
func (l *Logger) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.file.Sync()
}

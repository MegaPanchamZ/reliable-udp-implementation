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
	EventSend       EventType = "snd"
	EventReceive    EventType = "rcv"
	EventDrop       EventType = "drop"
	EventCorrupt    EventType = "crpt"
	EventDuplicate  EventType = "dup"
	EventTimeout    EventType = "timeout"
	EventRetransmit EventType = "rexmt"
)

// Logger provides thread-safe logging to a file
type Logger struct {
	file  *os.File
	mu    sync.Mutex
	start time.Time
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

// LogEvent logs a general event with timestamp
// Format: <timestamp> <event_type> <details>
func (l *Logger) LogEvent(eventType EventType, details string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := l.getTimestamp()
	line := fmt.Sprintf("%.2f %s %s\n", timestamp, eventType, details)
	l.file.WriteString(line)
}

// LogSegment logs a segment send/receive event
// Format: <timestamp> <event_type> <seq_num> <ack_num> <flags> <payload_size>
func (l *Logger) LogSegment(eventType EventType, seqNum, ackNum uint16, flags string, payloadSize int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := l.getTimestamp()
	line := fmt.Sprintf("%.2f %s %d %d %s %d\n",
		timestamp, eventType, seqNum, ackNum, flags, payloadSize)
	l.file.WriteString(line)
}

// LogDrop logs a packet drop event
// Format: <timestamp> drop <seq_num>
func (l *Logger) LogDrop(seqNum uint16) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := l.getTimestamp()
	line := fmt.Sprintf("%.2f drop %d\n", timestamp, seqNum)
	l.file.WriteString(line)
}

// LogCorrupt logs a packet corruption event
// Format: <timestamp> crpt <seq_num>
func (l *Logger) LogCorrupt(seqNum uint16) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := l.getTimestamp()
	line := fmt.Sprintf("%.2f crpt %d\n", timestamp, seqNum)
	l.file.WriteString(line)
}

// LogTimeout logs a timeout event
// Format: <timestamp> timeout <seq_num>
func (l *Logger) LogTimeout(seqNum uint16) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := l.getTimestamp()
	line := fmt.Sprintf("%.2f timeout %d\n", timestamp, seqNum)
	l.file.WriteString(line)
}

// LogRetransmit logs a retransmission event
// Format: <timestamp> rexmt <seq_num>
func (l *Logger) LogRetransmit(seqNum uint16) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := l.getTimestamp()
	line := fmt.Sprintf("%.2f rexmt %d\n", timestamp, seqNum)
	l.file.WriteString(line)
}

// LogStateChange logs a state transition
// Format: <timestamp> state <old_state> -> <new_state>
func (l *Logger) LogStateChange(oldState, newState string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := l.getTimestamp()
	line := fmt.Sprintf("%.2f state %s -> %s\n", timestamp, oldState, newState)
	l.file.WriteString(line)
}

// LogMessage logs a general message
func (l *Logger) LogMessage(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := l.getTimestamp()
	line := fmt.Sprintf("%.2f %s\n", timestamp, message)
	l.file.WriteString(line)
}

// Flush ensures all buffered data is written to disk
func (l *Logger) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.file.Sync()
}

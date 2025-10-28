package urp

import (
	"encoding/binary"
	"fmt"
)

// URPSegment represents a URP protocol segment
// Header format (6 bytes):
// - Sequence Number (2 bytes)
// - ACK Number (2 bytes)
// - Flags (1 byte): ACK, SYN, FIN
// - Checksum (1 byte, stored as uint16 but using 8 bits effectively)
type URPSegment struct {
	SeqNum   uint16
	AckNum   uint16
	Flags    uint8
	Checksum uint8
	Payload  []byte
}

// NewSegment creates a new URP segment
func NewSegment(seqNum, ackNum uint16, flags uint8, payload []byte) *URPSegment {
	seg := &URPSegment{
		SeqNum:  seqNum,
		AckNum:  ackNum,
		Flags:   flags,
		Payload: payload,
	}
	seg.ComputeAndSetChecksum()
	return seg
}

// Pack converts the segment into a byte slice for transmission
func (s *URPSegment) Pack() []byte {
	totalLen := HeaderSize + len(s.Payload)
	data := make([]byte, totalLen)

	// Pack header
	binary.BigEndian.PutUint16(data[0:2], s.SeqNum)
	binary.BigEndian.PutUint16(data[2:4], s.AckNum)
	data[4] = s.Flags
	data[5] = s.Checksum // Store 8-bit checksum

	// Pack payload
	if len(s.Payload) > 0 {
		copy(data[HeaderSize:], s.Payload)
	}

	return data
}

// Unpack parses a byte slice into a URPSegment
func Unpack(data []byte) (*URPSegment, error) {
	if len(data) < HeaderSize {
		return nil, fmt.Errorf("data too short: got %d bytes, need at least %d", len(data), HeaderSize)
	}

	seg := &URPSegment{
		SeqNum:   binary.BigEndian.Uint16(data[0:2]),
		AckNum:   binary.BigEndian.Uint16(data[2:4]),
		Flags:    data[4],
		Checksum: data[5],
	}

	// Extract payload if present
	if len(data) > HeaderSize {
		seg.Payload = make([]byte, len(data)-HeaderSize)
		copy(seg.Payload, data[HeaderSize:])
	}

	return seg, nil
}

// ComputeAndSetChecksum calculates and sets the checksum for this segment
func (s *URPSegment) ComputeAndSetChecksum() {
	// Create a temporary buffer with header fields and payload
	tempData := make([]byte, 5+len(s.Payload)) // Exclude checksum field itself
	binary.BigEndian.PutUint16(tempData[0:2], s.SeqNum)
	binary.BigEndian.PutUint16(tempData[2:4], s.AckNum)
	tempData[4] = s.Flags
	if len(s.Payload) > 0 {
		copy(tempData[5:], s.Payload)
	}

	s.Checksum = ComputeChecksum(tempData)
}

// IsValid checks if the segment's checksum is valid
func (s *URPSegment) IsValid() bool {
	// Recompute checksum and compare
	tempData := make([]byte, 5+len(s.Payload))
	binary.BigEndian.PutUint16(tempData[0:2], s.SeqNum)
	binary.BigEndian.PutUint16(tempData[2:4], s.AckNum)
	tempData[4] = s.Flags
	if len(s.Payload) > 0 {
		copy(tempData[5:], s.Payload)
	}

	computed := ComputeChecksum(tempData)
	return computed == s.Checksum
}

// HasFlag checks if a specific flag is set
func (s *URPSegment) HasFlag(flag uint8) bool {
	return (s.Flags & flag) != 0
}

// String returns a human-readable representation of the segment
func (s *URPSegment) String() string {
	flags := ""
	if s.HasFlag(FlagACK) {
		flags += "ACK "
	}
	if s.HasFlag(FlagSYN) {
		flags += "SYN "
	}
	if s.HasFlag(FlagFIN) {
		flags += "FIN "
	}
	if flags == "" {
		flags = "NONE"
	}

	return fmt.Sprintf("Seg[seq=%d ack=%d flags=%s len=%d]",
		s.SeqNum, s.AckNum, flags, len(s.Payload))
}

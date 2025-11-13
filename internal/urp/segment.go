package urp

import (
	"encoding/binary"
	"fmt"
)

// URPSegment represents a URP protocol segment
// Header format (6 bytes) per spec:
// - Sequence Number (2 bytes) - bytes 0-1
// - Reserved (13 bits) + Flags (3 bits) - bytes 2-3
// - Error Detection/Checksum (2 bytes) - bytes 4-5
//
// The sequence number field serves dual purpose:
// - For DATA/SYN/FIN: indicates sequence number
// - For ACK: indicates acknowledgment number (next expected byte)
type URPSegment struct {
	SeqNum   uint16 // Sequence number OR acknowledgment number (for ACK segments)
	Flags    uint16 // Upper 3 bits are flags (ACK, SYN, FIN), lower 13 bits reserved (must be 0)
	Checksum uint16 // 16-bit checksum
	Payload  []byte
}

// NewSegment creates a new URP segment
// seqNum: sequence number for DATA/SYN/FIN, or ack number for ACK segments
// flags: control flags (ACK, SYN, FIN)
func NewSegment(seqNum uint16, flags uint8, payload []byte) *URPSegment {
	seg := &URPSegment{
		SeqNum:  seqNum,
		Flags:   uint16(flags) << 13, // Shift flags to upper 3 bits
		Payload: payload,
	}
	seg.ComputeAndSetChecksum()
	return seg
}

// Pack converts the segment into a byte slice for transmission
func (s *URPSegment) Pack() []byte {
	totalLen := HeaderSize + len(s.Payload)
	data := make([]byte, totalLen)

	// Pack header per spec:
	// Bytes 0-1: Sequence Number
	binary.BigEndian.PutUint16(data[0:2], s.SeqNum)
	// Bytes 2-3: Reserved (13 bits) + Flags (3 bits)
	binary.BigEndian.PutUint16(data[2:4], s.Flags)
	// Bytes 4-5: Checksum (16 bits)
	binary.BigEndian.PutUint16(data[4:6], s.Checksum)

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
		Flags:    binary.BigEndian.Uint16(data[2:4]),
		Checksum: binary.BigEndian.Uint16(data[4:6]),
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
	// Create a temporary buffer with header fields (excluding checksum) and payload
	tempData := make([]byte, 4+len(s.Payload)) // SeqNum (2) + Flags (2) + Payload
	binary.BigEndian.PutUint16(tempData[0:2], s.SeqNum)
	binary.BigEndian.PutUint16(tempData[2:4], s.Flags)
	if len(s.Payload) > 0 {
		copy(tempData[4:], s.Payload)
	}

	s.Checksum = ComputeChecksum(tempData)
}

// IsValid checks if the segment's checksum is valid
func (s *URPSegment) IsValid() bool {
	// Recompute checksum and compare
	tempData := make([]byte, 4+len(s.Payload))
	binary.BigEndian.PutUint16(tempData[0:2], s.SeqNum)
	binary.BigEndian.PutUint16(tempData[2:4], s.Flags)
	if len(s.Payload) > 0 {
		copy(tempData[4:], s.Payload)
	}

	computed := ComputeChecksum(tempData)
	return computed == s.Checksum
}

// HasFlag checks if a specific flag is set
// Flags are in the upper 3 bits of the Flags field
func (s *URPSegment) HasFlag(flag uint8) bool {
	flagMask := uint16(flag) << 13
	return (s.Flags & flagMask) != 0
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
		flags = "DATA"
	}

	return fmt.Sprintf("Seg[seq=%d flags=%s len=%d]",
		s.SeqNum, flags, len(s.Payload))
}

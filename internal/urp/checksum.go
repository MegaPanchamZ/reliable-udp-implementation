package urp

// ComputeChecksum calculates a 16-bit Internet checksum (RFC 1071)
// This is the same algorithm used by TCP/UDP
func ComputeChecksum(data []byte) uint16 {
	var sum uint32

	// Process data in 16-bit words
	for i := 0; i < len(data)-1; i += 2 {
		word := uint32(data[i])<<8 | uint32(data[i+1])
		sum += word
	}

	// Add the last byte if data length is odd
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}

	// Fold 32-bit sum to 16 bits
	for sum>>16 > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}

	// Return one's complement
	return uint16(^sum)
}

// ValidateChecksum verifies the checksum of a segment
// Returns true if the checksum is valid
func ValidateChecksum(data []byte, checksum uint16) bool {
	computed := ComputeChecksum(data)
	return computed == checksum
}

// CorruptData intentionally corrupts a single byte in the data
// Used by the PLC module to simulate corruption
func CorruptData(data []byte) {
	if len(data) > 0 {
		// Flip a bit in the first byte
		data[0] ^= 0x01
	}
}

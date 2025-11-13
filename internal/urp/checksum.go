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

// CorruptData intentionally corrupts a single byte in the data per spec:
// "select a random byte in the segment (excluding the first four header bytes)
// and flip a single bit within that byte"
func CorruptData(data []byte, rng interface{ Intn(int) int }) {
	// Must exclude first 4 header bytes (SeqNum + first 2 bytes of Flags field)
	if len(data) <= 4 {
		return // Cannot corrupt if only header
	}

	// Select random byte from position 4 onwards
	byteIndex := 4 + rng.Intn(len(data)-4)

	// Select random bit (0-7) to flip
	bitIndex := rng.Intn(8)

	// Flip the selected bit
	data[byteIndex] ^= (1 << bitIndex)
}

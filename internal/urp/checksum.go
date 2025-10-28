package urp

// ComputeChecksum calculates an 8-bit checksum for the given data
// Using a simple sum modulo 256
func ComputeChecksum(data []byte) uint8 {
	var sum uint32

	for _, b := range data {
		sum += uint32(b)
	}

	// Return lower 8 bits
	return uint8(sum & 0xFF)
} // ValidateChecksum verifies the checksum of a segment
// Returns true if the checksum is valid
func ValidateChecksum(data []byte, checksum uint8) bool {
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

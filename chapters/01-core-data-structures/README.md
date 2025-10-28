# Chapter 1: Core Data Structures

Welcome to the workshop! Every great project starts with a solid foundation. For our Reliable UDP protocol, that foundation is the data structure that will carry our messages. This chapter will guide you from the basic problems of unreliable networks to creating the core data modules for our project.

---

### 1. Remember & Understand (The "Why")

**The Problem:** Imagine you're sending a novel to a friend using postcards. If you just write the story and send them, you'll face chaos:
-   **Out-of-Order:** Postcards arrive in the wrong sequence.
-   **Corrupted:** The ink gets smudged and becomes unreadable.

UDP is just like this. It sends data but offers no guarantees about order or integrity.

**The Theory:** To solve this, we need to create a standardized format for our data, a **Segment**. This is a digital envelope with two parts: a **Header** (instructions for the receiver) and a **Payload** (the actual data).

To fix the problems, our header needs two key fields:
-   A **Sequence Number (`SeqNum`)** to solve the ordering problem.
-   A **Checksum** to detect if the data has been corrupted.

---

### 2. Apply (The "How")

**Design Sketch (Workshop on Paper):**

Let's think about how to represent this in code. We need a way to handle the `Segment` itself and a separate utility for the checksum math.

1.  **The `Segment` Module:** How would you structure a class or a set of functions to handle this? You'll need to:
    *   Store the `SeqNum`, `Checksum`, and `Payload`.
    *   Have a function to `pack()` this data into a single stream of bytes to be sent over the network.
    *   Have a function to `unpack()` a stream of bytes from the network back into your structure.

2.  **The `Checksum` Module:** This should be a simple helper.
    *   It needs a function to `compute()` the checksum from a chunk of data.
    *   It needs another function to `validate()` that a chunk of data matches its checksum.

Sketch this out on paper. What would the function signatures look like? What data types would you use?

---

### 3. Analyze (The "Trade-offs")

**Puzzle / Critical Thinking:**

Our checksum algorithm will be a simple 8-bit sum of all the bytes in the packet (modulo 256). This is fast, but is it perfectly reliable?

Consider two different payloads: `[0x01, 0x02]` and `[0x02, 0x01]`.
1.  What is the 8-bit sum checksum for both?
2.  Now consider the payload `[0x01, 0x02, 0xFF]`. What is its checksum? What about `[0x01, 0x03, 0xFE]`?
3.  What does this tell you about the limitations of this simple checksum? What kind of data corruption might it *fail* to detect?

This is a classic trade-off between performance and robustness.

---

### 4. Evaluate (The "Justification")

**Design Rationale:**

The reference implementation uses an 8-bit checksum for simplicity and educational clarity. As you discovered in the puzzle, this method is not perfect. It can fail to detect errors where byte order is swapped or where multiple errors cancel each other out.

Professional protocols like TCP and UDP use a more robust algorithm called a **16-bit one's complement sum**. It's more computationally intensive but catches a much wider range of errors, including swapped bytes.

For this project, our simple checksum is sufficient to learn the *principle* of error detection. We are evaluating that educational clarity is more important here than cryptographic-level error detection.

---

### 5. Create (The "Implementation")

**Coding Workshop:**

It's time to write the code. We will create two files in the `internal/urp/` directory to handle these core data structures.

**File 1: `internal/urp/segment.go`**
This file will define our `URPSegment` struct and handle the packing/unpacking logic. Notice how it uses Go's `binary` package to handle converting the 16-bit numbers to network byte order (Big Endian).

```go
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
// - Checksum (1 byte)
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
	data[5] = s.Checksum

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
```

**File 2: `internal/urp/checksum.go`**
This file contains our simple error detection logic.

```go
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
}

// ValidateChecksum is a helper that can be used if you have the data and checksum separately
func ValidateChecksum(data []byte, checksum uint8) bool {
	computed := ComputeChecksum(data)
	return computed == checksum
}
```

You have now built the foundational data structures for the entire protocol!

**[Next Chapter: Connection Management](./../02-connection-management/README.md)**

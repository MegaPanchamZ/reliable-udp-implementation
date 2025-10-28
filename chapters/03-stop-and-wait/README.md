# Chapter 3: Stop-and-Wait Transfer

With a reliable connection handshake, we can now send data. We will start with the simplest possible reliable transfer algorithm: **Stop-and-Wait**. It's not the most efficient, but it's the perfect stepping stone to understanding more complex concepts.

---

### 1. Remember & Understand (The "Why")

**The Problem:**
We have an `ESTABLISHED` connection. The sender needs to send a file, broken into many segments. How does it know that each segment has arrived safely before sending the next one?

**The Theory: Stop-and-Wait**
The algorithm is as simple as its name suggests:
1.  **Sender:** Sends **one** data segment.
2.  **Sender:** **Stops** and **waits**.
3.  **Receiver:** Receives the segment and sends an **Acknowledgment (ACK)**.
4.  **Sender:** Receives the ACK, knows the segment arrived safely, and is now allowed to send the *next* segment.

This process repeats for every single segment. It's like a very cautious conversation where you wait for a "got it" after every sentence.

---

### 2. Apply (The "How")

**Design Sketch (Workshop on Paper):**

Let's add this logic to our existing `Sender` and `Receiver` modules.

1.  **Sender's Logic (`ESTABLISHED` state):**
    *   How will the sender keep track of which segment it's waiting for an ACK for? (Hint: `send_base`).
    *   After sending a segment with `SeqNum=X`, what `AckNum` should it expect back from the receiver?
    *   The sender's main loop will now involve: reading a chunk from the file, sending it, and waiting for a specific ACK before reading the next chunk.

2.  **Receiver's Logic (`ESTABLISHED` state):**
    *   How will the receiver know if a data segment is the one it's supposed to get next? (Hint: `expected_seq_num`).
    *   When it receives the correct segment, what two things must it do? (One involves the file, the other involves sending a message back).
    *   What should the `AckNum` be in its response?

---

### 3. Analyze (The "Trade-offs")

**Puzzle / Critical Thinking:**

Stop-and-Wait is very simple, but it has a major performance problem related to something called **Bandwidth-Delay Product (BDP)**.

Imagine a network connection between Earth and a Mars rover. The time for a radio signal to travel one way (the delay) is about 10 minutes. Let's say our connection has a high bandwidth of 1 Mbps.

1.  Sender sends a 1000-byte packet. This takes about 8 milliseconds to transmit.
2.  The packet travels to Mars (10 minutes).
3.  The rover receives it and sends an ACK.
4.  The ACK travels back to Earth (10 minutes).

How much time did the sender spend **actively sending** versus **waiting**? What percentage of the available 1 Mbps bandwidth is actually being used? This inefficiency is the primary trade-off of Stop-and-Wait.

---

### 4. Evaluate (The "Justification")

**Design Rationale:**

As the puzzle demonstrates, Stop-and-Wait is incredibly inefficient on networks with high latency. The sender spends most of its time idle, waiting for an ACK to travel across the network. The "pipe" of the network is mostly empty.

So why are we implementing it?
1.  **Simplicity:** It's the easiest way to start. The logic for managing the connection is trivial—we only ever need to worry about one packet at a time.
2.  **Foundation:** It forces us to build the core logic of sending data and processing ACKs. The more advanced "Sliding Window" protocol we'll build in Chapter 5 is a direct evolution of Stop-and-Wait, designed specifically to solve this efficiency problem by keeping the pipe full.

We are choosing to build this "inefficient" protocol first because it's the best way to learn the fundamentals before adding complexity.

---

### 5. Create (The "Implementation")

**Coding Workshop:**

Let's add the Stop-and-Wait logic to the `ESTABLISHED` state in our `Sender` and `Receiver`.

**File 1: `internal/protocol/sender.go` (Updated `SendFile` logic)**
We'll modify the main loop to send one segment and wait for its specific ACK.

```go
// Inside the Sender struct...

func (s *Sender) sendData(data []byte) error {
	offset := 0
	for offset < len(data) {
		// Determine the size of the next chunk
		end := offset + urp.MSS
		if end > len(data) {
			end = len(data)
		}
		payload := data[offset:end]

		// Create and send the segment
		seg := urp.NewSegment(s.nextSeqNum, 0, 0, payload)
		s.sendSegment(seg, true)

		// STOP and WAIT for the ACK
		ack, err := s.waitForAck(s.rto)
		if err != nil {
			// (In the next chapter, we'll add retransmission here)
			return fmt.Errorf("timeout waiting for ACK for segment %d", s.nextSeqNum)
		}

		// Check if the ACK is the one we expect
		if ack.AckNum == s.nextSeqNum+uint16(len(payload)) {
			// It is! Update our state and move to the next chunk.
			s.nextSeqNum = ack.AckNum
			offset = end
		} else {
			// This is a duplicate or old ACK, ignore it for now.
			// The timeout will handle retransmission if needed.
		}
	}
	return nil
}
```

**File 2: `internal/protocol/receiver.go` (Updated `Listen` logic)**
We'll add the data handling part to our `ESTABLISHED` state.

```go
// Inside the Receiver struct's Listen loop, in the ESTABLISHED case...

case urp.StateESTABLISHED:
    if seg.HasFlag(urp.FlagFIN) {
        // (FIN logic remains the same)
        // ...
        return nil
    }

    // Handle incoming data
    if len(seg.Payload) > 0 {
        // Is this the segment we are expecting?
        if seg.SeqNum == r.expectedSeq {
            // Yes. Write data to file.
            _, err := r.outputFile.Write(seg.Payload)
            if err != nil {
                return err
            }
            
            // Update what we expect next.
            r.expectedSeq += uint16(len(seg.Payload))
        }
        
        // Always send a cumulative ACK for the next byte we expect.
        // If we received a duplicate, this re-sends the needed ACK.
        ackSeg := urp.NewSegment(0, r.expectedSeq, urp.FlagACK, nil)
        r.sendAck(ackSeg)
    }
```

You have now implemented a basic, but functional, reliable transport protocol! It can transfer an entire file correctly over a perfect network. The next step is to make it robust enough to handle an imperfect one.

**[Previous Chapter: Connection Management](./../02-connection-management/README.md)** | **[Next Chapter: Robust Transfer](./../04-robust-transfer/README.md)**

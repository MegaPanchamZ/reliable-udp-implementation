# Chapter 2: Connection Management

With our core data structures in place, we can now address a fundamental question: how does a connection start and end? We can't just start sending data into the void. We need a polite "handshake" to begin the conversation and a graceful "goodbye" to end it.

---

### 1. Remember & Understand (The "Why")

**The Problem:**
1.  **Cold Start:** If the Sender just sends a data packet, the Receiver has no idea a new connection is starting. It might be a stray packet from a previous, crashed session. The Receiver needs to be explicitly told, "A new sender wants to talk to you."
2.  **Ambiguous End:** When the Sender is done, how does the Receiver know the file is complete? If the final packet is lost, the Receiver might wait forever for more data.

**The Theory: State Machines & Flags**
To solve this, we introduce the concept of a **connection state**. Both the Sender and Receiver will exist in a specific state (e.g., `LISTEN`, `ESTABLISHED`, `CLOSED`). The connection moves from one state to another based on special messages.

We'll use the `Flags` field in our `URPSegment` header to send these special messages:
-   **`SYN` (Synchronize):** A request to start a new connection.
-   **`FIN` (Finish):** A notification that one side is done sending data.
-   **`ACK` (Acknowledge):** A flag to indicate that the `AckNum` field is valid.

The process of using these flags to establish a connection is called a **three-way handshake**.

---

### 2. Apply (The "How")

**Design Sketch (Workshop on Paper):**

Let's design the logic for the handshake.

1.  **Sender's Role:**
    *   What's the first thing the Sender should send? (Hint: It needs to `SYN`chronize).
    *   After sending the `SYN`, what state should it be in? What is it waiting for?
    *   What does it need to receive from the Receiver to confirm the connection is open?

2.  **Receiver's Role:**
    *   What is the Receiver's initial state?
    *   What message should it be listening for?
    *   When it receives the `SYN`, what should it send back? (Hint: It needs to `ACK` the `SYN`).
    *   After sending its response, what state should it be in?

Sketch this out as a timeline. Draw the Sender on the left, the Receiver on the right, and use arrows to show the sequence of `SYN` and `ACK` messages.

---

### 3. Analyze (The "Trade-offs")

**Puzzle / Critical Thinking:**

The three-way handshake is `SYN` -> `SYN-ACK` -> `ACK`. Why is the third message (the final `ACK` from the sender) necessary? Why not just have a two-way handshake (`SYN` -> `SYN-ACK`)?

Consider a scenario where a `SYN` packet from an old connection is delayed in the network and arrives very late. What could go wrong if the Receiver acted on this old `SYN` and the handshake was only two steps? How does the third step prevent this problem?

---

### 4. Evaluate (The "Justification")

**Design Rationale:**

The third step of the handshake is crucial for preventing zombie connections. As you analyzed in the puzzle, a delayed `SYN` from a past connection could trick a server into allocating resources for a connection that the client has long since abandoned.

The third `ACK` serves as a live confirmation from the client: "Yes, I am the one who just sent that `SYN`, and I am ready to proceed." Without it, the server has no way of knowing if the `SYN` it received is fresh or a relic. The three-way handshake ensures that both sides have agreed to start a new connection *in the present moment*.

Similarly, for teardown, a `FIN` handshake ensures both sides know that the data stream is complete and can release their resources gracefully.

---

### 5. Create (The "Implementation")

**Coding Workshop:**

We will now create the "brain" of our protocol: the `Sender` and `Receiver` logic modules. These will contain the state machines that handle the connection lifecycle.

**File 1: `internal/protocol/sender.go` (Initial Version)**
This file will manage the sender's state. For now, we'll just implement the handshake and teardown logic.

```go
package protocol

import (
	"urp-go/internal/logger"
	"urp-go/internal/plc"
	"urp-go/internal/urp"
	"fmt"
	"net"
	"time"
)

// Sender implements the sender-side URP protocol
type Sender struct {
	// ... (other fields like conn, logger, plc)
	state      int
	sendBase   uint16
	nextSeqNum uint16
	// ...
}

// SendFile is the main entry point
func (s *Sender) SendFile(data []byte) error {
	// 1. Send SYN and wait for SYN-ACK
	if err := s.establishConnection(); err != nil {
		return err
	}

	// (Data transfer logic will go here in the next chapter)
	fmt.Println("Connection established!")

	// 2. Send FIN and wait for FIN-ACK
	if err := s.closeConnection(); err != nil {
		return err
	}

	fmt.Println("Connection closed.")
	return nil
}

func (s *Sender) establishConnection() error {
	s.changeState(urp.StateSYNSENT)
	synSeg := urp.NewSegment(s.nextSeqNum, 0, urp.FlagSYN, nil)
	
	// Simple retry logic for the handshake
	for i := 0; i < 3; i++ { // Try 3 times
		s.sendSegment(synSeg, true)
		
		ack, err := s.waitForAck(2 * time.Second) // 2-second timeout
		if err == nil && ack.HasFlag(urp.FlagSYN) && ack.HasFlag(urp.FlagACK) {
			s.sendBase = ack.AckNum
			s.nextSeqNum = ack.AckNum
			s.changeState(urp.StateESTABLISHED)
			return nil
		}
	}
	return fmt.Errorf("handshake failed: no SYN-ACK received")
}

func (s *Sender) closeConnection() error {
	s.changeState(urp.StateFINSENT)
	finSeg := urp.NewSegment(s.nextSeqNum, 0, urp.FlagFIN, nil)

	// Simple retry logic for teardown
	for i := 0; i < 3; i++ {
		s.sendSegment(finSeg, true)

		ack, err := s.waitForAck(2 * time.Second)
		if err == nil && ack.HasFlag(urp.FlagACK) {
			s.changeState(urp.StateCLOSED)
			return nil
		}
	}
	return fmt.Errorf("teardown failed: no FIN-ACK received")
}

// (Helper functions like sendSegment, waitForAck, changeState would also be here)
```

**File 2: `internal/protocol/receiver.go` (Initial Version)**
This file will manage the receiver's state.

```go
package protocol

import (
	"urp-go/internal/logger"
	"urp-go/internal/plc"
	"urp-go/internal/urp"
	"fmt"
	"net"
	"time"
)

// Receiver implements the receiver-side URP protocol
type Receiver struct {
	// ... (fields like conn, logger, plc)
	state       int
	expectedSeq uint16
	// ...
}

// Listen is the main loop
func (r *Receiver) Listen() error {
	for {
		seg, _, err := r.receiveSegment() // Simplified receive logic
		if err != nil {
			continue
		}

		switch r.state {
		case urp.StateLISTEN:
			if seg.HasFlag(urp.FlagSYN) {
				r.changeState(urp.StateESTABLISHED) // Simplified state change for now
				r.expectedSeq = seg.SeqNum + 1
				
				// Send SYN-ACK
				ackSeg := urp.NewSegment(0, r.expectedSeq, urp.FlagSYN|urp.FlagACK, nil)
				r.sendAck(ackSeg)
			}
		case urp.StateESTABLISHED:
			if seg.HasFlag(urp.FlagFIN) {
				r.changeState(urp.StateTIMEWAIT)
				
				// Send FIN-ACK
				ackSeg := urp.NewSegment(0, seg.SeqNum+1, urp.FlagACK, nil)
				r.sendAck(ackSeg)

				// In a real implementation, we'd wait here before closing.
				// For now, we'll just exit.
				return nil 
			}
			// (Data handling logic will go here)
		}
	}
}

// (Helper functions like receiveSegment, sendAck, changeState would also be here)
```

You have now built the logic to open and close a connection reliably. The next step is to send data within that connection.

**[Previous Chapter: Core Data Structures](./../01-core-data-structures/README.md)** | **[Next Chapter: Stop-and-Wait Transfer](./../03-stop-and-wait/README.md)**

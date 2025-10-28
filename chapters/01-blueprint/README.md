# Chapter 1: The Blueprint - Designing a Reliable Message

Welcome to the workshop! Our first task is to design the blueprint for our **URP Segment**. Think of this as a standardized digital envelope that will carry our data. We can't just send raw data; we need a structure to overcome the unreliability of UDP.

Let's build this structure step-by-step, solving one problem at a time.

---

### Problem 1: Out-of-Order Messages

Imagine you send three packets: `A`, `B`, and `C`. Because of the chaotic nature of the internet, they might arrive as `C`, `A`, `B`. The receiver has no idea how to reassemble them correctly.

**How can the receiver know the correct order?**

<details>
  <summary><strong>Solution: A Sequence Number</strong></summary>
  
  We add a **Sequence Number (`SeqNum`)** to our header. This is a simple counter.

  - Packet `A` gets `SeqNum = 1`
  - Packet `B` gets `SeqNum = 2`
  - Packet `C` gets `SeqNum = 3`

  Now, even if they arrive as `C(3)`, `A(1)`, `B(2)`, the receiver can use the `SeqNum` to put them back in the right order. This is the first piece of our blueprint.

  ```
  +------------------+------------------+
  |      Header      |     Payload      |
  |------------------|                  |
  |   SeqNum: 2      | (Data chunk B)   |
  +------------------+------------------+
  ```
</details>

---

### Problem 2: Lost Messages & Acknowledgments

We send packets 1, 2, and 3. Packet 2 gets lost in the void. The sender thinks the job is done, and the receiver is stuck forever waiting for a packet that will never arrive.

**How can the sender know if a packet was lost? And how can the receiver tell the sender what it has received?**

<details>
  <summary><strong>Solution: An Acknowledgment Number</strong></summary>
  
  The receiver needs to send messages back to the sender. These are called **Acknowledgments (ACKs)**. We'll add an **Acknowledgment Number (`AckNum`)** field to our header.

  The `AckNum` is the receiver's way of saying: "I have received all data up to this sequence number, and I am now expecting the *next* one."

  - **Sender sends `SeqNum=1`**.
  - **Receiver gets `SeqNum=1`** and sends back a packet with **`AckNum=2`**.
  - The sender sees `AckNum=2` and knows that packet 1 arrived safely. It can now send packet 2.

  If the sender doesn't receive an `AckNum=2` after a certain amount of time (a timeout), it assumes packet 1 was lost and sends it again. This process is called **retransmission**.

  ```
  +------------------+------------------+
  |      Header      |     Payload      |
  |------------------|                  |
  |   SeqNum: 2      | (Data chunk B)   |
  |   AckNum: 1      |                  |
  +------------------+------------------+
  ```
</details>

---

### Problem 3: Corrupted Data (Smudged Ink)

A packet arrives, but its data was flipped during transit due to a hardware error. The `SeqNum` and `AckNum` might be correct, but the payload is garbage. The receiver has no way of knowing the data is corrupt.

**How can we detect if the data has been tampered with or corrupted?**

<details>
  <summary><strong>Solution: A Checksum</strong></summary>
  
  We'll add a **Checksum** field. A checksum is the result of a mathematical function run over the packet's data.

  1.  **The Sender:** Before sending, it calculates the checksum of the packet and puts the result in the `Checksum` field.
  2.  **The Receiver:** When a packet arrives, it runs the *exact same* checksum function on the received data.
  3.  **Verification:** It compares its calculated checksum with the value in the `Checksum` field. If they don't match, the packet is corrupt and must be discarded.

  For our project, we'll use a simple 8-bit sum, but more complex algorithms (like CRC32) are common in real-world protocols.

  ```
  +------------------+------------------+
  |      Header      |     Payload      |
  |------------------|                  |
  |   SeqNum: 2      | (Data chunk B)   |
  |   AckNum: 1      |                  |
  |   Checksum: 184  |                  |
  +------------------+------------------+
  ```
</details>

---

### Problem 4: Special Messages (Starting & Ending the Conversation)

How do we start a connection? The sender can't just send data out of the blue. It needs to "synchronize" with the receiver first. Similarly, how do we end the connection gracefully?

**How can we send control messages that aren't part of the data stream?**

<details>
  <summary><strong>Solution: Flags</strong></summary>
  
  We'll reserve a single byte in our header for **Flags**. Each bit in this byte represents an "on/off" switch for a specific control message.

  - **`SYN` (Synchronize):** The "let's start a connection" flag. Used in the initial handshake.
  - **`ACK` (Acknowledge):** Indicates that the `AckNum` in this packet is valid. Most packets will have this on.
  - **`FIN` (Finish):** The "I'm done sending data" flag. Used to tear down the connection.

  These flags allow us to manage the connection state without mixing control messages into our data payload.
</details>

---

### Our Final Blueprint: The URP Segment

By solving these problems, we have designed the complete header for our URP Segment.

| Field | Size (Bytes) | Description |
|---|---|---|
| **SeqNum** | 2 | Sequence number of this packet. |
| **AckNum** | 2 | Sequence number the sender is expecting next. |
| **Flags** | 1 | Control flags (SYN, ACK, FIN). |
| **Checksum**| 1 | For detecting data corruption. |
| **Payload** | Up to 1000 | The actual file data. |

This 6-byte header is the heart of our protocol. It contains all the "magic instructions" needed to transform unreliable UDP into a reliable data stream.

### Your Turn: Think Ahead

1.  **The Handshake:** How would you use the `SYN` and `ACK` flags to create a reliable connection startup, similar to a polite conversation? (e.g., "I'd like to talk." -> "Okay, I'm listening." -> "Great, here's the first message.")
2.  **Data vs. ACKs:** Can a single packet carry both data (payload) and an acknowledgment (`AckNum`)? Why would this be efficient?

**[Next Chapter: The Messenger](./../02-messenger/README.md)** - Now that we have our blueprint, let's build the sender that uses it.
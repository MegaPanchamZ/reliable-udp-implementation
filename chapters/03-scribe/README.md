# Chapter 3: The Scribe - Building the Receiver

The Receiver's primary duty is to be a meticulous scribe: reassemble the sender's data perfectly and write it to the destination file. Its logic is simpler than the sender's, but it must be precise.

Let's build the Receiver's logic by solving its main challenges.

---

### Problem 1: Out-of-Order Delivery

The network is chaotic. Packets arrive out of order. The application layer (the file we are writing) requires the data in the correct, sequential order. The receiver cannot simply write data to the file as it arrives.

**How can the receiver ensure it delivers an ordered stream of data to the application?**

<details>
  <summary><strong>Solution: Buffering and `expected_seq_num`</strong></summary>
  
  The receiver will maintain two key components:

  1.  **`expected_seq_num`**: A variable that holds the sequence number of the *next* piece of data it needs to write to the file. It's the receiver's most important state.
  2.  **A Receive Buffer**: A temporary storage area (like a dictionary or hash map) for packets that arrive out of order.

  **The Logic:**
  - When a packet arrives, the receiver checks if its sequence number matches `expected_seq_num`.
  - **If it matches:** The data is written to the file. The receiver then checks its buffer to see if the next sequential packets have already arrived. If so, they are also written to the file, and `expected_seq_num` is updated accordingly.
  - **If it does NOT match:** The packet is placed in the receive buffer, waiting for the missing piece to arrive.

</details>

---

### Problem 2: Communicating Receipt

The receiver needs to inform the sender what it has received, which drives the sender's sliding window and retransmission logic.

**What information should the receiver send back, and when?**

<details>
  <summary><strong>Solution: Cumulative Acknowledgments (ACKs)</strong></summary>
  
  The receiver sends an `ACK` segment back to the sender for every data packet it receives. The key is that these ACKs are **cumulative**.

  **The Rule:** The `AckNum` field in the receiver's response is *always* set to `expected_seq_num`.

  This single number elegantly tells the sender everything it needs to know:
  - "I have received all data perfectly up to `expected_seq_num - 1`."
  - "I am now waiting for the packet with sequence number `expected_seq_num`."

  This design is what enables both the sliding window (by advancing `send_base`) and fast retransmit (by generating duplicate ACKs when out-of-order packets arrive).

</details>

---

### Problem 3: Handling Corrupt and Duplicate Packets

Packets can be corrupted by the network or duplicated due to retransmissions. The receiver must handle both cases gracefully.

**What should the receiver do with a corrupt or duplicate packet?**

<details>
  <summary><strong>Solution: Validate, then Acknowledge or Discard</strong></summary>
  
  The receiver follows a strict sequence of checks for every incoming packet:

  1.  **Validate Checksum:** First, it calculates the checksum of the received packet. If it doesn't match the checksum in the header, the packet is corrupt. **It is silently discarded.** The sender's RTO will handle the retransmission. No ACK is sent.
  2.  **Check for Duplicates:** If the checksum is valid, the receiver checks the sequence number. If it's a number it has already received and written to the file (i.e., `SeqNum < expected_seq_num`), the data is a duplicate. The data is discarded, but an ACK is **re-sent** with the current `expected_seq_num`. This is critical, as the sender's previous ACK may have been the packet that was lost.

</details>

---

### The Receiver's State Machine

The receiver's state machine is simpler than the sender's:

1.  **`CLOSED`**: The initial state.
2.  **`LISTEN`**: The receiver is waiting for the initial `SYN` packet from the sender.
3.  **`ESTABLISHED`**: The main state. The receiver is accepting data packets, managing its buffer, and sending ACKs.
4.  **`TIME_WAIT`**: After receiving a `FIN` from the sender, the receiver sends a final `ACK` and enters this state. It waits for a short period (e.g., 2 seconds) to ensure any lingering packets die out before closing the connection completely. This prevents old duplicate packets from a closed connection from interfering with a new connection on the same ports.

With this logic, the Receiver acts as the perfect counterpart to the Sender, enabling a fully reliable data transfer.

**[Previous Chapter: The Messenger](./../02-messenger/README.md)** | **[Next Chapter: The Gremlin](./../04-gremlin/README.md)**
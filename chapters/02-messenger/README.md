# Chapter 2: The Messenger - Building the Sender

Now that we have our `URPSegment` blueprint, we can build the Sender. The Sender's job is complex: it must send data, ensure it arrives, and manage the flow of information without overwhelming the network.

Let's build the Sender's logic by tackling its core challenges one by one.

---

### Problem 1: Flow Control - How to Avoid Flooding the Network?

If the sender sends 1,000 packets all at once, it could easily overwhelm the receiver's buffer or congest the network path. This is like trying to force a novel through a mail slot instead of sending it page by page.

**How can the sender limit the number of packets it sends without waiting for an ACK for every single one?**

<details>
  <summary><strong>Solution: The Sliding Window</strong></summary>
  
  We'll implement a **Sliding Window**. This is a conceptual "window" that defines how much data the sender can have "in-flight" (sent but not yet acknowledged).

  - **Window Size (`max_win`):** A limit, in bytes, on the amount of unacknowledged data.
  - **`send_base`:** The sequence number of the oldest unacknowledged packet. This is the left edge of the window.
  - **`next_seq_num`:** The sequence number of the next new packet to be sent. This is the right edge of the window.

  **The Rule:** The sender can only send new packets if `next_seq_num - send_base < max_win`.

  When an ACK arrives, it "slides" the window forward by updating `send_base`. This allows the sender to transmit more data. This approach is far more efficient than Stop-and-Wait, where the window size is always just one packet.

</details>

---

### Problem 2: Packet Loss - What If an ACK Never Comes?

The sender sends a packet, but the packet (or its corresponding ACK) is lost. Without a mechanism to handle this, the sender would wait forever, and the connection would stall.

**How does the sender detect and recover from a lost packet?**

<details>
  <summary><strong>Solution: The Retransmission Timeout (RTO)</strong></summary>
  
  The sender will use a **timer**.

  1.  **Start Timer:** When the sender sends a packet, it starts a timer for the *oldest unacknowledged segment* (`send_base`).
  2.  **Wait for ACK:** If the ACK for `send_base` arrives, the window slides, and the timer is restarted for the new `send_base` (if there's still unacknowledged data).
  3.  **Timeout:** If the timer expires before the ACK arrives, the sender assumes the packet was lost. It **retransmits** the segment at `send_base` and restarts the timer.

  **The Goldilocks Problem:** Choosing the right RTO value is critical.
  - **Too short:** Unnecessary retransmissions, causing network congestion.
  - **Too long:** The protocol feels slow and sluggish when packets are dropped.

</details>

---

### Problem 3: Inefficiency - Waiting for a Timeout is Slow

Imagine packets 10, 11, 12, and 13 are sent. Packet 11 is lost, but 12 and 13 arrive at the receiver. The sender has to wait for the full RTO to expire before it retransmits packet 11, even though there's strong evidence that a specific packet is missing.

**Is there a faster way to detect a single lost packet in a stream?**

<details>
  <summary><strong>Solution: Fast Retransmit</strong></summary>
  
  We can use duplicate ACKs as a clue. The receiver always ACKs the next in-order sequence number it's expecting.

  - Receiver gets packet 10. It sends `ACK=11`.
  - Packet 11 is lost.
  - Receiver gets packet 12. It's out of order. It discards 12 and sends another `ACK=11`. (This is a **duplicate ACK**).
  - Receiver gets packet 13. It's also out of order. It discards 13 and sends a third `ACK=11`.

  **The Rule:** When the sender receives **three duplicate ACKs** for the same sequence number, it's a very strong signal that the packet immediately following that ACK was lost. The sender can then **immediately retransmit** the missing packet without waiting for the RTO timer to expire. This dramatically improves performance on networks with occasional packet loss.

</details>

---

### The Sender's State Machine

Combining these solutions, the sender operates as a **state machine**. Here's a simplified view:

1.  **`CLOSED`**: The initial state.
2.  **`SYN_SENT`**: After sending the initial `SYN` to start a connection, wait for a `SYN-ACK`.
3.  **`ESTABLISHED`**: The main state for data transfer. Here, it manages the sliding window, RTO timer, and fast retransmit logic.
4.  **`FIN_WAIT`**: After sending a `FIN` to close the connection, wait for the final `ACK`.

By implementing this logic, our Sender becomes a robust messenger capable of handling the chaos of an unreliable network.

**[Previous Chapter: The Blueprint](./../01-blueprint/README.md)** | **[Next Chapter: The Scribe](./../03-scribe/README.md)**
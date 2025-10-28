# Chapter 2: The Messenger - The Sender's Tale

The Sender has the hardest job. It must be a meticulous bookkeeper, a patient watchman, and a swift courier all at once. Its logic is built around two key concepts: the **Sliding Window** and the **Retransmission Timer**.

### The Sliding Window: A Conductor's Baton

The Sender can't just send all 100 pages of the novel at once. The network (and the receiver) would be overwhelmed. It needs to control the flow of data. The sliding window algorithm is the solution.

Imagine the Sender has a "window" of a certain size, say 4 pages. This means it can send pages 1, 2, 3, and 4 without waiting for an acknowledgment. But it cannot send page 5 until it knows page 1 has been received safely.

-   `sendBase`: The start of the window. The oldest page sent but not yet acknowledged.
-   `nextSeqNum`: The end of the window. The next new page to send.

When the receiver sends an `ACK` for page 1, the window "slides" forward. The `sendBase` is now 2, and the Sender is allowed to send page 5.

### "Draw It Out" Challenge #3

Let's visualize this. Assume a window size of 5.
1.  Draw the initial state. The Sender has sent pages 1, 2, 3, 4, 5. The window covers these pages. `sendBase` is 1, `nextSeqNum` is 6. The window is full.
2.  Now, an `ACK` arrives for page 3. **Important:** Our ACKs are cumulative. An ACK for page 3 means "I have received everything up to and including page 3."
3.  Draw the "after" state. Where does the window slide to? What is the new `sendBase`? What new pages can the Sender now send?

### The Retransmission Timer: A Test of Patience

What if page 3 was lost? The Sender would wait forever for an `ACK` that will never come. To prevent this, the Sender starts a timer every time it sends data. If the timer goes off before an `ACK` arrives, it assumes the packet was lost and **retransmits** it. This is the Retransmission Timeout (RTO).

### Puzzle: The Goldilocks Timer

What problems would occur if your RTO value was:
1.  **Too short?** (e.g., you retransmit before the original packet even has a chance to arrive).
2.  **Too long?** (e.g., you wait for ages before realizing a packet was lost).

<details>
  <summary>Click to reveal the answer</summary>
  
  1.  **Too short:** You would flood the network with unnecessary duplicate packets, causing congestion and wasting bandwidth. It's like a nervous person who keeps asking "Did you get my text?" every five seconds.
  2.  **Too long:** The connection would feel sluggish and slow. Users would experience long stalls whenever a packet is dropped.
</details>

### Fast Retransmit: An Optimization

Waiting for a timeout is slow. There's a clever trick to speed things up. If the Sender receives the *same* `ACK` multiple times (usually 3), it's a huge clue.

Imagine the Sender sent pages 10, 11, 12, 13. Page 11 gets lost, but 12 and 13 arrive. The Receiver will send:
-   `ACK 11` (when it receives page 10)
-   `ACK 11` (when it receives page 12, it's still waiting for 11)
-   `ACK 11` (when it receives page 13, it's *still* waiting for 11)

When the Sender sees three duplicate `ACK 11`s, it doesn't wait for the timer. It immediately knows page 11 is the culprit and retransmits it. This is **Fast Retransmit**.

**[Previous Chapter: The Blueprint](./../01-blueprint/README.md)** | **[Next Chapter: The Scribe](./../03-scribe/README.md)**

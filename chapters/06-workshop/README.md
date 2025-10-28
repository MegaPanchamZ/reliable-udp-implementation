# Chapter 6: The Workshop - A Phased Implementation Guide

This guide provides a logical, step-by-step order to build your URP implementation. Following these phases will help you build from simple to complex, with a testable milestone at each step. You will rely heavily on your log files at every stage.

---

### Phase 1: The Skeleton (Argument & Socket Setup)

Your goal here is to just make the two programs talk. Forget the protocol for a second.

1.  **Argument Parsing:** Write the code for both `Sender` and `Receiver` to parse all their respective command-line arguments.
2.  **Socket Setup:**
    *   **Receiver:** Create a UDP socket and bind it to the `receiver_port`. Make it print a message when it receives *any* data.
    *   **Sender:** Create a UDP socket. Don't bind it (or bind to `sender_port`). Write code to send a simple "Hello" string to the `receiver_port`.
3.  **Test:** Run the Receiver, then the Sender. If the Receiver prints "Hello," you know your basic socket and argument logic is correct.

---

### Phase 2: Core Data Structures (Segment & Logging)

Now, build the tools you'll need. Don't write any protocol logic yet.

1.  **`URPSegment` Class/Module:** Create a class or module to handle your segment.
    *   Implement `pack()`: A function that takes `seq_num`, flags, and data, and builds the 6-byte header plus payload into a single `bytes` object.
    *   Implement `unpack()`: A function that takes raw `bytes` from the socket and parses them into a `URPSegment` object.
2.  **`ErrorDetection` Module:** Create a helper module with two functions:
    *   `compute_checksum(bytes)`: Implements your chosen error-detection scheme.
    *   `validate_checksum(bytes)`: Returns `True` or `False`.
3.  **`Logger` Class:** Create a logging utility to write to `sender_log.txt` and `receiver_log.txt`. This is **critical for debugging**. Get the timestamping and formatting right. You will be using this in *every* step from now on.

---

### Phase 3: Stop-and-Wait on a *Reliable* Channel

This is your "Minimum Viable Product." Set `max_win = 1000` and all four loss/corruption probabilities to `0`.

1.  **Implement Handshake:**
    *   **Sender:** Send a `SYN` segment. Wait for an `ACK`.
    *   **Receiver:** In a `LISTEN` state, wait for a `SYN`. When it arrives, send the `ACK` and move to `ESTABLISHED`.
2.  **Implement Data Transfer (Stop-and-Wait):**
    *   **Sender:** Read 1000 bytes (MSS) from the file. Send it as one `DATA` segment. Wait for the correct `ACK` before sending the next segment.
    *   **Receiver:** If in `ESTABLISHED`, wait for a `DATA` segment. If the sequence number is the one you expect, write the data to the file and send a cumulative `ACK`.
3.  **Implement Teardown:**
    *   **Sender:** After the last file segment is ACKed, send a `FIN` segment. Wait for the `ACK`.
    *   **Receiver:** When you get a `FIN`, send an `ACK`, move to `TIME_WAIT`. Wait 2 seconds, then close.

**Test:** Your file should transfer perfectly. Your logs should be clean.

---

### Phase 4: Stop-and-Wait on an *Unreliable* Channel

Now, make your Phase 3 code robust. Keep `max_win = 1000`.

1.  **Implement the PLC Module:** Add the Packet Loss and Corruption logic to your `Sender`. Make *all* outgoing and incoming segments pass through it.
2.  **Implement Timers (Sender):** This is the hard part.
    *   Add the `rto` timer.
    *   Start the timer when you send `SYN`, `DATA`, or `FIN`.
    *   If the timer expires, retransmit the segment (`SYN`, `DATA`, or `FIN`) and restart the timer.
    *   If you receive the correct `ACK`, stop the timer.
3.  **Implement Checksum & Duplicates (Receiver):**
    *   Use your `validate_checksum()` function on *every* segment. If it fails, **silently discard the packet** (do not send an `ACK`). The Sender's timeout will handle it.
    *   Handle duplicate packets. If you receive a `SYN` or `DATA` segment you've already processed, discard the data but **re-send the ACK**. This is because your previous ACK might have been lost.

**Test:** Now, run with `flp`, `rlp`, `fcp`, and `rcp` > 0. Your file should *still* transfer correctly, and your logs will show `drp`, `cor`, and retransmissions.

---

### Phase 5: Sliding Window & Fast Retransmit

This is the final upgrade from Stop-and-Wait to full URP.

1.  **Sliding Window (Sender):**
    *   Upgrade your Sender's logic. It now needs a buffer for all unacknowledged segments.
    *   It should send new segments as long as the window (`next_seq_num - send_base`) is less than `max_win`.
    *   The timer logic changes: It *only* runs for the oldest unacknowledged segment (`send_base`). When an ACK arrives that slides the window, restart the timer *if* there are still un-ACKed segments in flight.
2.  **Buffering & Cumulative ACKs (Receiver):**
    *   Upgrade your Receiver's logic. It needs a receive buffer.
    *   Your ACK logic *must* be cumulative. The `ACK` number should always be the *next byte you expect in order*.
    *   If you receive an out-of-order segment, buffer it. Send a *duplicate ACK* for the byte you're still waiting for.
    *   When the missing segment arrives, write all contiguous, in-order data (from your buffer) to the file and send a new, updated cumulative `ACK`.
3.  **Fast Retransmit (Sender):**
    *   Implement the `dup_ack_count`.
    *   If you receive a duplicate `ACK`, increment the count.
    *   If `dup_ack_count == 3`, immediately retransmit the oldest unacknowledged segment and reset the count.

---

### Phase 6: Finalization

1.  **Statistics:** Implement the logic to count all the final statistics and append them to the log files.
2.  **Final Test:** Test *everything* in a consistent environment to ensure robustness.

This order builds from simple to complex, and each phase is a testable milestone. Good luck!

**[Previous Chapter: The Grand Assembly](./../05-assembly/README.md)**
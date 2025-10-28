# Chapter 5: The Grand Assembly - Your Pseudocode Guide

Now it's time to bring everything together. Here is a language-agnostic pseudocode skeleton for your URP implementation. This is your guide to writing the actual code.

### Sender Main Logic

```
function send_file(file_data):
  // Chapter 1: Perform the three-way handshake
  establish_connection()

  send_base = 0
  next_seq_num = 0
  window_buffer = new map() // Stores sent but un-acked segments

  while send_base < len(file_data):
    // Chapter 2: Fill the sliding window
    while window_is_not_full() and next_seq_num < len(file_data):
      segment = create_segment(file_data, next_seq_num)
      send(segment) // This goes through the Gremlin (PLC)
      window_buffer.add(segment)
      start_timer_if_not_running()
      next_seq_num += len(segment.payload)

    // Wait for an ACK or for the timer to expire
    event = wait_for_ack_or_timeout()

    if event is ack:
      // Update send_base based on the cumulative ACK
      new_send_base = ack.ack_num
      // Remove acknowledged segments from buffer
      window_buffer.remove_all_before(new_send_base)
      send_base = new_send_base
      // Handle fast retransmit logic
      check_for_duplicate_acks(ack)
      // Restart timer if there are still un-acked segments
      restart_timer_if_needed()
    else if event is timeout:
      // Retransmit the oldest un-acked segment
      retransmit_segment(window_buffer.get(send_base))
      restart_timer()

  // Send the FIN to close the connection
  close_connection()
```

### Receiver Main Logic

```
function receive_file():
  // Chapter 1: Listen for and respond to the initial SYN
  wait_for_connection()

  expected_seq_num = 0
  out_of_order_buffer = new map() // Stores future segments

  while connection_is_not_closed:
    segment = receive_packet() // This comes from the Gremlin (PLC)

    // Chapter 1: Check for smudges
    if checksum_is_invalid(segment):
      continue // Silently discard the corrupted segment

    // Chapter 3: Handle the data
    if segment.seq_num == expected_seq_num:
      write_to_file(segment.payload)
      expected_seq_num += len(segment.payload)

      // Check buffer for any segments that can now be delivered
      while out_of_order_buffer.has(expected_seq_num):
        buffered_segment = out_of_order_buffer.get(expected_seq_num)
        write_to_file(buffered_segment.payload)
        expected_seq_num += len(buffered_segment.payload)
        out_of_order_buffer.delete(buffered_segment.seq_num)

    else if segment.seq_num > expected_seq_num:
      // Buffer future segments
      out_of_order_buffer.set(segment.seq_num, segment)

    // Always send an ACK for what you're currently expecting
    send_ack(expected_seq_num)
```

**[Previous Chapter: The Gremlin](./../04-gremlin/README.md)** | **[Next Chapter: The Workshop](./../06-workshop/README.md)**

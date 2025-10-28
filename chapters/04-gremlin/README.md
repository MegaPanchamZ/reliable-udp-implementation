# Chapter 4: The Gremlin - Simulating Chaos

How do we know if our protocol actually works? We can't just hope for a bad network day. We must create our own chaos. We will build a **Packet Loss and Corruption (PLC) Module**—a gremlin that lives between our Sender/Receiver and the actual network.

The PLC's only job is to randomly drop and corrupt segments based on probabilities we give it.

### Puzzle: Taming the Gremlin

You want to simulate a network with a 25% packet loss rate. In your PLC module, which gets to inspect every outgoing packet, how would you implement this? Describe the logic in simple terms or pseudocode.

<details>
  <summary>Click to reveal the answer</summary>
  
  For every packet that passes through the PLC, generate a random number between 0 and 1.
  ```
  function should_i_drop_this_packet():
    loss_probability = 0.25
    random_value = generate_random(0.0, 1.0)

    if random_value < loss_probability:
      return TRUE  // Drop the packet
    else:
      return FALSE // Let the packet pass
  ```
</details>

By creating this gremlin, we can rigorously test our protocol under harsh but controlled conditions, ensuring it's robust enough for the real world. The `run_tests.sh` script in the root of this repository uses this very logic to test our implementation.

**[Previous Chapter: The Scribe](./../03-scribe/README.md)** | **[Next Chapter: The Grand Assembly](./../05-assembly/README.md)**

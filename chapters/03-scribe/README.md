# Chapter 3: The Scribe - The Receiver's Duty

The Receiver's job is to be an impeccable scribe. It gathers the incoming postcards, puts them in the correct order, and writes the final, perfect copy of the novel.

### Order from Chaos: Buffering

The Receiver has one main rule: **deliver data to the application in perfect order.**

What happens if page 5 arrives before page 4? The Receiver can't just deliver page 5. It must hold onto it and wait for the missing page 4. This holding area is called a **buffer**.

### "Draw It Out" Challenge #4

The Receiver is expecting page 20 (`expectedSeqNum = 20`). It receives segments in the following order: `SeqNum=22`, `SeqNum=21`, `SeqNum=20`.

1.  Draw a box representing the Receiver's buffer.
2.  Show the state of the buffer after `SeqNum=22` arrives. What `ACK` number does the Receiver send back? (Hint: It's still waiting for 20!).
3.  Update the drawing for when `SeqNum=21` arrives. What `ACK` is sent?
4.  Finally, show what happens when the long-awaited `SeqNum=20` arrives. What happens to the buffer? What is the new `expectedSeqNum`?

This exercise demonstrates the core logic of the Receiver:
1.  If the received segment is the one I expect, accept it. Then, check my buffer to see if I can accept any more now.
2.  If the received segment is from the future, put it in the buffer.
3.  No matter what, always send an `ACK` for the sequence number I am currently waiting for.

### The Power of the ACK

The `ACK` is the Receiver's only voice. It's a powerful tool that simultaneously tells the Sender two things:
1.  "I have successfully received everything up to this point."
2.  "This is the next piece of data I am expecting."

This elegant, dual-purpose design is what makes the sliding window and fast retransmit algorithms possible.

**[Previous Chapter: The Messenger](./../02-messenger/README.md)** | **[Next Chapter: The Gremlin](./../04-gremlin/README.md)**

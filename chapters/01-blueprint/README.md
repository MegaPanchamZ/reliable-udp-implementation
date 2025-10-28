# Chapter 1: The Blueprint - The Anatomy of a Magic Scroll

Before we can send our postcard novel, we need a system. We can't just scribble on a piece of paper and hope for the best. We need a standardized format, a magic scroll that carries our message and the instructions to understand it. In networking, this is called a **Segment**.

Our URP Segment is the atomic unit of our protocol. It's a small data structure with two parts: the **Header** (the instructions) and the **Payload** (the actual page of our novel).

### Puzzle: The Essential Instructions

You are designing the header for our segment. To solve the postcard problem (lost, out-of-order, and smudged pages), what is the absolute minimum information you need to include in the header?

<details>
  <summary>Click to reveal the answer</summary>
  
  1.  **To fix out-of-order pages:** We need a **Sequence Number** (`SeqNum`). This is like writing "Page 1 of 100," "Page 2 of 100," etc., on each postcard.
  2.  **To fix lost pages:** The receiver needs a way to tell the sender, "I'm missing page 2!" This is done with an **Acknowledgment Number** (`AckNum`). The receiver uses it to say, "I've received everything up to page 1, and I'm now waiting for page 2."
  3.  **To fix smudged pages:** We need a way to detect if the data has been corrupted. A **Checksum** is a "magic number" calculated from the data. If the receiver calculates a different checksum, it knows the data is smudged and can discard it.
  4.  **For special messages:** We need a way to start and end the conversation gracefully. We'll use **Flags**, like `SYN` (synchronize, "I want to start talking!") and `FIN` (finish, "I'm done talking!").
</details>

### "Draw It Out" Challenge #1

On your paper, draw a rectangle representing our URP Segment. Divide it into a Header and a Payload. Now, based on the puzzle solution, sketch out the fields inside the Header. Don't worry about the exact size in bytes yet, just the concepts. Your drawing should look something like this:

```
+--------------------------------------------------+
|                 URP Segment                      |
+----------------------+---------------------------+
|        Header        |          Payload          |
|                      |                           |
|  - Sequence Number   |   (A chunk of the file)   |
|  - Ack Number        |                           |
|  - Flags (SYN, FIN)  |                           |
|  - Checksum          |                           |
+----------------------+---------------------------+
```

This blueprint is the core of our entire protocol. Every packet we send will follow this structure.

### The Magic Handshake

We can't just start shouting data into the void. We need to establish a connection. In networking, this is often done with a **three-way handshake**. It's like a polite conversation:

1.  **Sender -> Receiver:** "Hello! I'd like to start sending. My first page is number 1." (This is a segment with the `SYN` flag set).
2.  **Receiver -> Sender:** "I hear you! I'm ready for page 1. Are you still there?" (A segment with `SYN` and `ACK` flags).
3.  **Sender -> Receiver:** "Yes, I'm here! Here comes the data." (A segment with the `ACK` flag).

The connection is now **ESTABLISHED**.

### "Draw It Out" Challenge #2

Draw a timeline for this handshake. Show the Sender on the left and the Receiver on the right. Use arrows to represent the three messages, and label them with the flags (`SYN`, `SYN|ACK`, `ACK`). This visual will be your guide for the first part of the sender and receiver logic.

**[Next Chapter: The Messenger](./../02-messenger/README.md)**

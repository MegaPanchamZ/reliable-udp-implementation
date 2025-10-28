# Chapter 1: Core Data Structures

Welcome to the workshop! Every great project starts with a solid foundation. For our Reliable UDP protocol, that foundation is the data structure that will carry our messages. This chapter will guide you from the basic problems of unreliable networks to creating the core data modules for our project.

---

### 1. Remember & Understand (The "Why")

**The Problem:** Imagine you're sending a novel to a friend using postcards. If you just write the story and send them, you'll face chaos:
-   **Out-of-Order:** Postcards arrive in the wrong sequence.
-   **Corrupted:** The ink gets smudged and becomes unreadable.

UDP is just like this. It sends data but offers no guarantees about order or integrity.

**The Theory:** To solve this, we need to create a standardized format for our data, a **Segment**. This is a digital envelope with two parts: a **Header** (instructions for the receiver) and a **Payload** (the actual data).

To fix the problems, our header needs two key fields:
-   A **Sequence Number (`SeqNum`)** to solve the ordering problem.
-   A **Checksum** to detect if the data has been corrupted.

---

### 2. Apply (The "How")

**Design Sketch (Workshop on Paper):**

Let's think about how to represent this in code. We need a way to handle the `Segment` itself and a separate utility for the checksum math.

1.  **The `Segment` Module:** How would you structure a class or a set of functions to handle this? You'll need to:
    *   Store the `SeqNum`, `Checksum`, and `Payload`.
    *   Have a function to `pack()` this data into a single stream of bytes to be sent over the network.
    *   Have a function to `unpack()` a stream of bytes from the network back into your structure.

2.  **The `Checksum` Module:** This should be a simple helper.
    *   It needs a function to `compute()` the checksum from a chunk of data.
    *   It needs another function to `validate()` that a chunk of data matches its checksum.

Sketch this out on paper. What would the function signatures look like? What data types would you use?

---

### 3. Analyze (The "Trade-offs")

**Puzzle / Critical Thinking:**

Our checksum algorithm will be a simple 8-bit sum of all the bytes in the packet (modulo 256). This is fast, but is it perfectly reliable?

Consider two different payloads: `[0x01, 0x02]` and `[0x02, 0x01]`.
1.  What is the 8-bit sum checksum for both?
2.  Now consider the payload `[0x01, 0x02, 0xFF]`. What is its checksum? What about `[0x01, 0x03, 0xFE]`?
3.  What does this tell you about the limitations of this simple checksum? What kind of data corruption might it *fail* to detect?

This is a classic trade-off between performance and robustness.

---

### 4. Evaluate (The "Justification")

**Design Rationale:**

The reference implementation uses an 8-bit checksum for simplicity and educational clarity. As you discovered in the puzzle, this method is not perfect. It can fail to detect errors where byte order is swapped or where multiple errors cancel each other out.

Professional protocols like TCP and UDP use a more robust algorithm called a **16-bit one's complement sum**. It's more computationally intensive but catches a much wider range of errors, including swapped bytes.

For this project, our simple checksum is sufficient to learn the *principle* of error detection. We are evaluating that educational clarity is more important here than cryptographic-level error detection.

---

### 5. Create (Your Implementation)

**Coding Workshop:**

It's time to write the code. Based on your design sketch, create the two core modules for your project.

1.  **Create your `Segment` module.** This class or struct should contain the header fields and payload. Implement `pack()` and `unpack()` methods.
    *   **Important:** When packing and unpacking multi-byte fields (like the 16-bit `SeqNum`), you must handle the **network byte order**, which is **big-endian**. Most languages provide library functions to ensure this conversion is done correctly.
2.  **Create your `ErrorDetection` module.** This should contain your `compute_checksum()` and `validate_checksum()` helper functions.

You have now built the foundational data structures for the entire protocol! For a complete reference, you can examine the Go implementation in the `internal/urp/` directory of this repository.

**[Next Chapter: Connection Management](./../02-connection-management/README.md)**

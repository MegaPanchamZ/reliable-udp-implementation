# The Unreliable Reliable Protocol: A Journey into Network Sorcery

Welcome, aspiring network sorcerer! This repository is not just a collection of code; it's an interactive book designed to guide you through the creation of your very own reliable transport protocol, URP.

You will learn to conquer the chaos of an unreliable network by crafting the data structures (spells) and algorithms (rituals) that guarantee data arrives perfectly, every time.

## The Adventure Ahead (Table of Contents)

This guide is structured into chapters. It's recommended to follow them in order, as each chapter builds upon the last.

*   **[Chapter 1: The Blueprint](./chapters/01-blueprint/README.md)**
    *   We design the core data structure of our protocol: the `URPSegment`.
    *   *Concepts: Headers, Payloads, Sequence Numbers, ACKs, Checksums, Flags.*

*   **[Chapter 2: The Messenger](./chapters/02-messenger/README.md)**
    *   We build the logic for the Sender, the brains of the operation.
    *   *Concepts: State Machines, Sliding Windows, Retransmission Timers, Fast Retransmit.*

*   **[Chapter 3: The Scribe](./chapters/03-scribe/README.md)**
    *   We create the Receiver, which patiently reassembles our data from chaos.
    *   *Concepts: Buffering, Cumulative ACKs, Flow Control.*

*   **[Chapter 4: The Gremlin](./chapters/04-gremlin/README.md)**
    *   We build a tool to simulate an unreliable network to test our protocol.
    *   *Concepts: Simulating Packet Loss, Corruption, and Reordering.*

*   **[Chapter 5: The Grand Assembly](./chapters/05-assembly/README.md)**
    *   We look at the complete pseudocode that brings all our concepts together.

*   **[Chapter 6: The Workshop](./chapters/06-workshop/README.md)**
    *   The practical guide to setting up your environment, with socket examples in multiple languages.

## Your Quest: Getting Started

You can approach this repository in two ways:

1.  **The Sorcerer's Apprentice (Reading & Exploration):** Read through the chapters in the `/chapters` directory. Examine the final source code in `/src` to see how the concepts are implemented.

2.  **The Master Crafter (Hands-On):** Follow the guides in each chapter and try to build the protocol yourself from scratch. Use the code in `/src` as a reference when you get stuck.

### Running the Tests

To prove our magic works, we have a gauntlet of tests. You can run them locally or in a controlled Docker environment.

**1. Local Testing (Requires Go)**

This is the best way to see the protocol in action and debug it.

```bash
# First, give the script execute permissions
chmod +x ./run_tests.sh

# Run the entire test suite
./run_tests.sh
```

**2. Isolated Testing (Requires Docker)**

This method builds the sender and receiver inside a container and runs a single test case to verify the environment. It's perfect for ensuring the protocol works consistently anywhere.

```bash
# This command will build the Docker image and run the sender/receiver services.
# The sender will transfer a file to the receiver with 10% packet loss.
docker-compose up --build
```
After the test completes, you can check the `data/` directory. If `data/test_file.txt` and `data/received_complete.txt` are identical, the test passed! You can bring down the environment with `Ctrl+C`.
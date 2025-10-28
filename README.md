# The Unreliable Reliable Protocol: A Journey into Network Sorcery

Welcome, aspiring network sorcerer! This repository is not just a collection of code; it's an **interactive educational resource** designed to guide you through the creation of your very own reliable transport protocol, **URP** (UDP Reliable Protocol).

You will learn to conquer the chaos of an unreliable network by crafting the data structures (spells) and algorithms (rituals) that guarantee data arrives perfectly, every time—just like TCP, but built by you!

## 🚀 Quick Start

**Want to see it work immediately?**

```bash
# Clone the repository
git clone https://github.com/MegaPanchamZ/reliable-udp-implementation.git
cd reliable-udp-implementation

# Build the executables
go build -o sender.exe ./cmd/sender
go build -o receiver.exe ./cmd/receiver

# Run quick validation tests (3 baseline tests)
./run_tests.ps1 -Quick    # Windows
# or
./run_tests.sh            # Linux/Mac
```

**For detailed build instructions, see [BUILDING.md](./BUILDING.md)**.

## 📚 The Adventure Ahead (Table of Contents)

This guide is structured into chapters. It's recommended to follow them in order, as each chapter builds upon the last.

*   **[Chapter 1: The Blueprint](./chapters/01-blueprint/README.md)**
    *   Design the core data structure: the `URPSegment`
    *   *Concepts: Headers, Payloads, Sequence Numbers, ACKs, Checksums, Flags*

*   **[Chapter 2: The Messenger](./chapters/02-messenger/README.md)**
    *   Build the Sender's state machine and sliding window logic
    *   *Concepts: State Machines, Sliding Windows, Retransmission Timers, Fast Retransmit*

*   **[Chapter 3: The Scribe](./chapters/03-scribe/README.md)**
    *   Create the Receiver that reassembles data from chaos
    *   *Concepts: Buffering, Cumulative ACKs, Flow Control, Out-of-Order Handling*

*   **[Chapter 4: The Gremlin](./chapters/04-gremlin/README.md)**
    *   Build a Packet Loss & Corruption (PLC) simulator
    *   *Concepts: Simulating Packet Loss, Corruption, Network Unreliability*

*   **[Chapter 5: The Grand Assembly](./chapters/05-assembly/README.md)**
    *   See the complete system pseudocode
    *   *Concepts: Integration, Connection Lifecycle, End-to-End Flow*

*   **[Chapter 6: The Workshop](./chapters/06-workshop/README.md)**
    *   Practical environment setup and socket programming examples
    *   *Multiple languages: Go, Python, C, JavaScript, and more*

## 🎯 Your Quest: Getting Started

You can approach this repository in multiple ways:

### 1. 📖 **The Sorcerer's Apprentice** (Learning Mode)
   - Read through the chapters in `/chapters` directory sequentially
   - Understand the theory and design decisions
   - Examine the reference implementation in `/internal` and `/cmd`
   - **Best for:** Understanding how reliable protocols work

### 2. 🛠️ **The Master Crafter** (Hands-On Mode)
   - Follow chapter guides and implement the protocol yourself
   - Use the reference code when stuck
   - Test your implementation against the test suite
   - **Best for:** Building practical networking skills

### 3. 🔬 **The Experimenter** (Research Mode)
   - Modify the existing protocol (change window size, timeouts, etc.)
   - Run tests to see the effects
   - Try the examples in different programming languages
   - **Best for:** Exploring protocol behavior and optimization

## 🧪 Testing Your Protocol

### Quick Validation
```bash
# Windows PowerShell
.\run_tests.ps1 -Quick

# Linux/Mac
./run_tests.sh --quick
```

### Full Test Suite (16 tests)
```bash
.\run_tests.ps1           # Windows
./run_tests.sh            # Linux/Mac
```

**Test Categories:**
- **Stop-and-Wait Protocol** (8 tests) - Window = 1×MSS
- **Sliding Window Protocol** (5 tests) - Window = 5×MSS  
- **Edge Cases** (3 tests) - Non-MSS segments, state verification

### Docker Testing
```bash
# Build and run in isolated container
docker-compose up --build

# Clean up
docker-compose down
```

After the test completes, check that `data/test_file.txt` and `data/received_complete.txt` are identical—that means it worked! ✅

## 📦 Repository Structure

```
UDPAssignmentFun/
├── cmd/                   # Executable programs
│   ├── sender/           # Sender application
│   └── receiver/         # Receiver application
├── internal/             # Core protocol implementation
│   ├── urp/             # Protocol definitions (segment, checksum, constants)
│   ├── protocol/        # State machines (sender & receiver FSMs)
│   ├── plc/             # Packet Loss & Corruption simulator
│   └── logger/          # Thread-safe event logging
├── chapters/            # Educational content (step-by-step guides)
├── examples/            # Code examples in 7+ languages
├── data/                # Test files and log outputs
├── BUILDING.md          # Detailed build and run instructions
├── CONTRIBUTING.md      # How to contribute to this project
└── README.md            # You are here!
```

## 🎓 What You'll Learn

By following this repository, you'll gain deep understanding of:

- **Network Programming** - UDP sockets, packet handling, addressing
- **Protocol Design** - Headers, checksums, sequence numbers, state machines
- **Reliability Mechanisms** - Retransmission, acknowledgments, duplicate detection
- **Flow Control** - Sliding windows, buffering, congestion management
- **Error Handling** - Corruption detection, timeout management, graceful failures
- **Concurrency** - Thread-safe logging, asynchronous I/O
- **Testing** - Simulation, edge cases, probabilistic scenarios

## 🤝 Contributing

We welcome contributions! Whether you're:
- Fixing typos
- Adding examples in new languages
- Improving documentation
- Creating visualizations
- Enhancing the protocol

See **[CONTRIBUTING.md](./CONTRIBUTING.md)** for guidelines.

## 📝 Protocol Specifications

**URP (UDP Reliable Protocol)** features:
- **Header:** 6 bytes (SeqNum:2, AckNum:2, Flags:1, Checksum:1)
- **MSS:** 1000 bytes maximum payload
- **Checksum:** 8-bit simple sum (Σbytes mod 256)
- **Flow Control:** Sliding window (configurable size)
- **Retransmission:** Fixed RTO (no exponential backoff)
- **Fast Retransmit:** Triggered on 3 duplicate ACKs
- **Connection:** 3-way handshake (SYN/SYN-ACK)
- **Teardown:** Graceful FIN handshake with TIME_WAIT

## 📖 Additional Resources

- **[BUILDING.md](./BUILDING.md)** - Complete build, run, and troubleshooting guide
- **[QUICKSTART.md](./QUICKSTART.md)** - Quick reference for common tasks
- **[TEST_SUITE.md](./TEST_SUITE.md)** - Detailed test documentation

## 🏆 Project Status

✅ **Fully Functional** - All core features implemented and tested
- Stop-and-wait mode: 100% reliable
- Sliding window: Handles packet loss/corruption with retransmission
- Comprehensive test suite with 16 validation scenarios
- Multi-language examples available

## 📄 License

This is an educational project for learning reliable transport protocol design. Feel free to use, modify, and learn from it!

## 🌟 Acknowledgments

Built as an educational resource to demystify how TCP-like protocols work under the hood. Inspired by classic networking courses and RFCs.

---

**Ready to begin your journey?** Start with [Chapter 1: The Blueprint](./chapters/01-blueprint/README.md) or jump straight to [BUILDING.md](./BUILDING.md) to run the code!
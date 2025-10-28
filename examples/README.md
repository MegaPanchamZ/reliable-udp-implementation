# UDP Programming Examples

This directory contains simple UDP sender/receiver examples in **7 programming languages**. These examples demonstrate the basics of UDP socket programming before diving into the full URP (UDP Reliable Protocol) implementation.

## Available Examples

| Language | Directory | Level | Description |
|----------|-----------|-------|-------------|
| **Go** | [go/](./go/) | Beginner | Simple examples with error handling |
| **Python** | [python/](./python/) | Beginner | Clean syntax, great for learning |
| **JavaScript** | [javascript/](./javascript/) | Beginner | Node.js with event-driven approach |
| **C** | [c/](./c/) | Intermediate | Low-level BSD sockets |
| **C++** | [cpp/](./cpp/) | Intermediate | Boost.Asio library |
| **Kotlin** | [kotlin/](./kotlin/) | Intermediate | JVM-based with null safety |
| **Dart** | [dart/](./dart/) | Intermediate | Async/await with streams |

## What These Examples Demonstrate

All examples implement a simple **echo protocol**:
1. **Sender** sends "Hello, URP!" to port 8080
2. **Receiver** listens on port 8080
3. **Receiver** echoes back with "Echo: " prefix
4. **Sender** prints the response: "Echo: Hello, URP!"

### Key Concepts Covered
- Creating UDP sockets
- Binding to a port (receiver)
- Sending datagrams
- Receiving datagrams
- Basic error handling
- Address/port management

## Quick Start

Each language directory has its own README with:
- Installation instructions
- Compilation steps (if needed)
- Running instructions
- Code highlights
- Learning resources

### General Pattern

**Terminal 1 - Start Receiver:**
```bash
cd examples/[language]
# Run receiver (see language-specific README)
```

**Terminal 2 - Run Sender:**
```bash
cd examples/[language]
# Run sender (see language-specific README)
```

## Comparing Languages

### Go
- Built-in UDP support (net package)
- Fast compilation
- Good error handling
- No external dependencies
- **No version requirement** - Any modern Go version works

### Python
- Easiest to learn
- Clean, readable syntax
- Built-in socket library
- Great for prototyping
- **Requires Python 3.9+** (3.8 and earlier are EOL)

### JavaScript (Node.js)
- Event-driven model
- Non-blocking I/O
- Web-friendly
- Large ecosystem
- **Requires Node.js 18+** (LTS recommended)

### C
- Direct system calls
- Maximum performance
- Low-level control
- **Cross-platform** - Works on Windows, Linux, macOS
- More complex error handling

### C++
- Boost.Asio - modern async I/O
- RAII for resource management
- Type safety
- Requires Boost library
- **Requires C++11 or later**

### Kotlin
- Null safety
- Java interoperability
- Concise syntax
- Requires JVM
- **Requires JDK 8+**

### Dart
- Modern async/await
- Stream-based
- Cross-platform
- Flutter integration
- **Requires Dart SDK 3.0+** (for null safety)

## Learning Path

1. **Start Simple** - Try Python or Go examples first
2. **Understand UDP** - Compare how each language handles sockets
3. **Add Complexity** - Try C for low-level understanding
4. **Explore Async** - Check JavaScript or Dart for event-driven models
5. **Study URP** - Move to `/cmd` and `/internal` for full protocol

## Important Notes

### These Are NOT Reliable
These examples use **raw UDP** without any reliability mechanisms:
- No acknowledgments
- No retransmission
- No ordering guarantees
- No error detection beyond OS-level
- No flow control

### Packets Can Be Lost
If you don't see a response, try running again. UDP is **best-effort delivery**.

### Port Conflicts
If port 8080 is in use, you'll get a "bind failed" error. Change the port in both sender and receiver.

## After These Examples

Once you understand basic UDP:

1. **Read the Chapters** - [/chapters](../chapters/) explains how to add reliability
2. **Study the Protocol** - [/internal/urp](../internal/urp/) shows segment structure
3. **Examine State Machines** - [/internal/protocol](../internal/protocol/) implements reliability
4. **Run the Full System** - [/cmd](../cmd/) contains complete sender/receiver with URP
5. **Run Tests** - [/run_tests.ps1](../run_tests.ps1) validates everything works

## Troubleshooting

### "Address already in use"
Another process is using port 8080. Kill it or change the port.

### "Connection refused" / No response
Make sure the receiver is running **before** starting the sender.

### Firewall blocks
Some firewalls block UDP. Allow port 8080 or disable firewall temporarily.

### Version Issues
- **Python:** Check version with `python --version` (need 3.9+)
- **Node.js:** Check version with `node --version` (need 18+)
- **Dart:** Check version with `dart --version` (need 3.0+)

### Windows: Missing compiler
- **C/C++:** Install MinGW-w64 or Visual Studio (with C++ tools)
- **Python/JS/Dart:** No compiler needed (interpreted)
- **Go/Kotlin:** Download from official sites

### Windows: C compilation errors
If you see errors about `sys/socket.h` or `unistd.h`:
- Make sure you're using the updated code (includes Windows support)
- Link with `-lws2_32` flag when using MinGW
- Use `cl` compiler with `ws2_32.lib` when using MSVC

## Additional Resources

- [UDP Protocol Overview](https://en.wikipedia.org/wiki/User_Datagram_Protocol)
- [TCP vs UDP](https://www.cloudflare.com/learning/ddos/glossary/user-datagram-protocol-udp/)
- [Socket Programming](https://beej.us/guide/bgnet/)
- [Main Project README](../README.md)

## Contributing

Want to add examples in another language? See [CONTRIBUTING.md](../CONTRIBUTING.md)!

Languages we'd love to see:
- Rust
- Swift
- Java
- C#
- Ruby
- PHP
- Elixir

---

**Ready to add reliability?** Check out the [main implementation](../cmd/) and [chapters](../chapters/) to learn how URP transforms these simple examples into a robust, reliable protocol!
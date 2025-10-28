# Go Examples - UDP Reliable Protocol

This directory contains simple Go examples to help you understand UDP socket programming basics before diving into the full URP implementation.

## Files

- **`sender.go`** - Minimal UDP sender (sends a message and waits for response)
- **`receiver.go`** - Minimal UDP receiver (listens and echoes back)

## Understanding the Basics

These examples demonstrate raw UDP communication **without** reliability. Notice that:
- No acknowledgments
- No retransmission
- No error checking
- Packets can be lost

The full URP protocol in `/cmd` and `/internal` adds all the reliability features on top of this basic UDP foundation.

## Running the Basic Examples

### Terminal 1: Start the Simple Receiver
```bash
cd examples/go
go run receiver.go
```

Output:
```
Simple UDP Receiver listening on :8080
Waiting for messages...
```

### Terminal 2: Run the Simple Sender
```bash
cd examples/go
go run sender.go
```

Output:
```
Sending message to server...
Received response: 'Echo: Hello, URP!'
```

## What's Happening?

1. **Receiver** creates a UDP socket and listens on port 8080
2. **Sender** creates a UDP socket and sends "Hello, URP!" to 127.0.0.1:8080
3. **Receiver** receives the message and sends back "Echo: ..." response
4. **Sender** receives the echo and prints it

## Try These Experiments

### Experiment 1: Packet Loss
Kill the receiver before the sender completes. Notice the sender waits indefinitely—there's no timeout or retry mechanism!

### Experiment 2: Multiple Messages
Modify sender.go to send multiple messages in a loop. Notice they might arrive out of order or not at all—there's no sequencing!

### Experiment 3: Large Messages
Try sending a message larger than the network's MTU (typically 1500 bytes). UDP will fragment it at the IP layer, but the application has no control or visibility into this.

## Key Differences from URP

| Feature | Basic UDP | URP Protocol |
|---------|-----------|--------------|
| **Reliability** | None - packets can be lost | Guaranteed delivery via ACKs and retransmission |
| **Ordering** | Not guaranteed | In-order delivery via sequence numbers |
| **Error Detection** | None in these examples | 8-bit checksum |
| **Flow Control** | None | Sliding window |
| **Connection** | Connectionless | Connection-oriented (3-way handshake) |
| **State** | Stateless | Stateful (5-state sender, 4-state receiver) |

## Next Steps

After understanding basic UDP:

1. **Read Chapter 1** - Learn how we structure data into segments
2. **Read Chapter 2** - See how the sender adds reliability
3. **Read Chapter 3** - Understand receiver buffering and ACKs
4. **Study `/cmd/sender/main.go`** - See the full sender implementation
5. **Study `/internal/protocol/sender.go`** - Deep dive into the state machine

## UDP Socket Fundamentals (Go)

### Creating a UDP Connection
```go
// Listen (server/receiver)
addr, _ := net.ResolveUDPAddr("udp", ":8080")
conn, _ := net.ListenUDP("udp", addr)

// Dial (client/sender)
serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
conn, _ := net.DialUDP("udp", nil, serverAddr)
```

### Sending Data
```go
data := []byte("Hello!")
conn.Write(data)
// or
conn.WriteToUDP(data, destinationAddr)
```

### Receiving Data
```go
buffer := make([]byte, 1024)
n, addr, _ := conn.ReadFromUDP(buffer)
message := buffer[:n]
```

### Important Notes
- UDP is **message-oriented** (not stream-oriented like TCP)
- Each `Write` sends one datagram
- Each `Read` receives one complete datagram
- Datagrams can arrive in any order
- Datagrams can be lost entirely
- No built-in acknowledgment or retransmission

## Learning Resources

- [Go net package documentation](https://pkg.go.dev/net)
- [UDP vs TCP comparison](https://www.cloudflare.com/learning/ddos/glossary/user-datagram-protocol-udp/)
- [Network programming in Go](https://pkg.go.dev/net)

## Questions?

- Check the main [README.md](../../README.md)
- Read [BUILDING.md](../../BUILDING.md) for troubleshooting
- Open an issue if you're stuck!

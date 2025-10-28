# Custom UDP Implementation (URP)

A custom reliable transport protocol built on top of UDP, implementing a sliding window mechanism similar to TCP. This project demonstrates fundamental concepts of reliable data transfer including connection management, flow control, and error recovery.

## Features

### Core Protocol Features
- **Custom Reliable Protocol (URP)**: Implements reliability on top of UDP
- **Sliding Window**: Efficient data transfer with configurable window size
- **Flow Control**: Prevents receiver buffer overflow
- **Connection Management**: SYN/ACK handshake and graceful FIN termination
- **Error Recovery**: 
  - Timeout-based retransmission
  - Fast retransmit on 3 duplicate ACKs
  - Checksum validation for data integrity

### State Machines
- **Sender (5 states)**: `CLOSED` → `SYN_SENT` → `ESTABLISHED` → `FIN_SENT` → `FIN_ACKED`
- **Receiver (4 states)**: `LISTEN` → `SYN_RCVD` → `ESTABLISHED` → `TIME_WAIT`

### Packet Loss Control (PLC) Module
Simulates network conditions with configurable parameters:
- Packet loss probability (forward and reverse paths)
- Packet corruption probability
- Packet delay and jitter
- Out-of-order delivery simulation

### Logging System
Comprehensive logging of all protocol events:
- Segment transmission/reception with sequence numbers
- State transitions
- Retransmissions and timeouts
- Statistical summaries

## Project Structure

```
.
├── cmd/
│   ├── sender/          # Sender application
│   └── receiver/        # Receiver application
├── internal/
│   ├── protocol/        # Core protocol implementation
│   │   ├── sender.go    # Sender logic with sliding window
│   │   └── receiver.go  # Receiver logic with buffering
│   ├── urp/            # URP segment and protocol definitions
│   │   ├── segment.go   # Segment structure and packing
│   │   ├── checksum.go  # Checksum calculation
│   │   └── const.go     # Protocol constants
│   ├── plc/            # Packet Loss Control simulation
│   │   └── module.go
│   └── logger/         # Event logging system
│       └── logger.go
├── data/               # Test files
├── docker-compose.yml  # Docker configuration for testing
└── Dockerfile
```

## Getting Started

### Prerequisites
- Go 1.22 or higher
- Docker (optional, for containerized testing)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/custom-udp-implementation.git
cd custom-udp-implementation
```

2. Build the applications:
```bash
go build -o sender.exe ./cmd/sender
go build -o receiver.exe ./cmd/receiver
```

### Usage

#### Basic File Transfer

1. Start the receiver:
```bash
./receiver.exe -port 8081 -sender_port 8080 -output received.txt
```

2. Start the sender:
```bash
./sender.exe -file data/test_file.txt -port 8080 -remote_host 127.0.0.1 -remote_port 8081
```

#### Advanced Options

**Sender Options:**
- `-file`: Input file path (required)
- `-port`: Local port (default: 8080)
- `-remote_host`: Receiver host (default: 127.0.0.1)
- `-remote_port`: Receiver port (default: 8081)
- `-window`: Maximum window size in bytes (default: 5000)
- `-rto`: Retransmission timeout in ms (default: 1000)
- `-log`: Log file path (default: sender_log.txt)

**Receiver Options:**
- `-port`: Local port (default: 8081)
- `-sender_port`: Sender port (default: 8080)
- `-output`: Output file path (default: received_file.txt)
- `-window`: Maximum window size in bytes (default: 5000)
- `-log`: Log file path (default: receiver_log.txt)

**PLC (Packet Loss Control) Options:**
Both sender and receiver support:
- `-loss_forward`: Forward path loss probability (0.0-1.0)
- `-loss_reverse`: Reverse path loss probability (0.0-1.0)
- `-corrupt`: Corruption probability (0.0-1.0)
- `-delay`: Base delay in ms (default: 0)
- `-jitter`: Max jitter variation in ms (default: 0)

#### Example with Packet Loss Simulation

```bash
# Receiver with 10% packet loss
./receiver.exe -port 8081 -loss_reverse 0.1

# Sender with 20% packet loss and 50ms delay
./sender.exe -file data/test_file.txt -loss_forward 0.2 -delay 50 -jitter 20
```

### Docker Usage

Run tests using Docker Compose:
```bash
docker-compose up
```

## Protocol Specification

### Segment Format

```
0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          Sequence Number      |       Acknowledgment Number   |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|     Flags     |              Checksum                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           Payload                             |
|                            (variable)                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

**Flags:**
- `SYN` (0x01): Synchronize connection
- `ACK` (0x02): Acknowledge
- `FIN` (0x04): Finish connection

### Protocol Constants
- **Header Size**: 7 bytes
- **MSS (Maximum Segment Size)**: 1000 bytes
- **Default Window Size**: 5000 bytes (5 segments)
- **Default RTO**: 1000ms
- **TIME_WAIT Duration**: 2000ms

## Testing

The project includes a PowerShell test script:

```powershell
.\run_tests.ps1
```

This script runs multiple test scenarios with different file sizes and network conditions.

## Implementation Details

### Sliding Window Protocol
- Maintains a window of unacknowledged segments
- Uses cumulative acknowledgments
- Implements Go-Back-N with selective repeat capability through buffering

### Reliability Mechanisms
1. **Sequence Numbers**: Track all data segments
2. **Checksums**: Detect corruption
3. **Timeouts**: Recover from packet loss
4. **Duplicate ACKs**: Fast retransmit after 3 duplicates
5. **Out-of-Order Buffering**: Receiver buffers future segments

### Performance Optimizations
- Non-blocking I/O with goroutines
- Channel-based communication for ACKs
- Efficient byte packing/unpacking
- Minimal memory allocation in hot paths

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

This project is provided as-is for educational purposes.

## Acknowledgments

This implementation is inspired by TCP/IP protocols and demonstrates fundamental concepts in reliable data transfer over unreliable networks.

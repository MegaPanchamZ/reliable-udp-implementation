# Building and Running the URP Protocol

This guide shows you how to build and test the UDP Reliable Protocol implementation.

## Prerequisites

- **Go 1.22 or later** - [Download Go](https://go.dev/dl/)
- **Git** (for cloning the repository)
- **Docker** (optional, for containerized testing)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/MegaPanchamZ/reliable-udp-implementation.git
cd reliable-udp-implementation
```

### 2. Build the Executables

```bash
# Build sender
go build -o sender.exe ./cmd/sender

# Build receiver
go build -o receiver.exe ./cmd/receiver
```

On Linux/Mac, omit the `.exe` extension:
```bash
go build -o sender ./cmd/sender
go build -o receiver ./cmd/receiver
```

### 3. Run a Simple Test

**Terminal 1 - Start the receiver:**
```bash
./receiver.exe --receiver_port 5000 --sender_port 6000 --file data/received.txt --max_win 4000
```

**Terminal 2 - Start the sender:**
```bash
./sender.exe --sender_port 6000 --receiver_host 127.0.0.1 --receiver_port 5000 --file data/test_file.txt --max_win 4000 --rto 200
```

### 4. Verify the Transfer

Check if the files match:
```bash
# On Windows (PowerShell)
Compare-Object (Get-Content data/test_file.txt) (Get-Content data/received.txt)

# On Linux/Mac
diff data/test_file.txt data/received.txt
```

No output means the files are identical! ✅

## Running the Full Test Suite

### PowerShell (Windows)
```powershell
# Run all 16 tests
.\run_tests.ps1

# Run only baseline tests (quick validation)
.\run_tests.ps1 -Quick
```

### Bash (Linux/Mac)
```bash
# Make script executable
chmod +x run_tests.sh

# Run all tests
./run_tests.sh
```

## Understanding the Test Results

The test suite validates:
- **Stop-and-Wait Protocol** (window = MSS): Tests 1.1-1.8
- **Sliding Window Protocol** (window = 5×MSS): Tests 2.1-2.5
- **Edge Cases**: Tests 3.1, 3.4, 3.5

**Expected Results:**
- Baseline tests (1.1, 2.1, 3.1) should always pass (100%)
- Tests with simulated packet loss/corruption may fail occasionally due to randomness
- A 60-80% overall pass rate is good for tests with network unreliability

Results are saved to `test_results.csv` with detailed metrics.

## Command-Line Parameters

### Sender Parameters

| Parameter | Description | Example | Required |
|-----------|-------------|---------|----------|
| `--sender_port` | Sender's UDP port | `6000` | Yes |
| `--receiver_host` | Receiver IP/hostname | `127.0.0.1` | Yes |
| `--receiver_port` | Receiver's UDP port | `5000` | Yes |
| `--file` | File to send | `data/test.txt` | Yes |
| `--max_win` | Window size (bytes) | `4000` | Yes |
| `--rto` | Retransmission timeout (ms) | `200` | Yes |
| `--flp` | Forward loss probability | `0.1` | No (default: 0.0) |
| `--rlp` | Reverse loss probability | `0.1` | No (default: 0.0) |
| `--fcp` | Forward corruption probability | `0.05` | No (default: 0.0) |
| `--rcp` | Reverse corruption probability | `0.05` | No (default: 0.0) |
| `--log` | Log file path | `data/sender_log.txt` | No |

### Receiver Parameters

| Parameter | Description | Example | Required |
|-----------|-------------|---------|----------|
| `--receiver_port` | Receiver's UDP port | `5000` | Yes |
| `--sender_port` | Sender's UDP port | `6000` | Yes |
| `--file` | Output file path | `data/received.txt` | Yes |
| `--max_win` | Window size (bytes) | `4000` | Yes |
| `--rlp` | Reverse loss probability | `0.1` | No (default: 0.0) |
| `--rcp` | Reverse corruption probability | `0.05` | No (default: 0.0) |
| `--log` | Log file path | `data/receiver_log.txt` | No |

## Testing with Packet Loss

Simulate unreliable networks:

```bash
# 10% packet loss on forward path (DATA segments)
./sender.exe --sender_port 6000 --receiver_host 127.0.0.1 --receiver_port 5000 \
  --file data/test_file.txt --max_win 4000 --rto 200 --flp 0.1

# 20% packet loss on reverse path (ACK segments)
./receiver.exe --receiver_port 5000 --sender_port 6000 \
  --file data/received.txt --max_win 4000 --rlp 0.2

# 5% corruption on both paths
./sender.exe ... --fcp 0.05 --rcp 0.05
```

## Viewing Logs

Logs show every segment sent/received:

```bash
# View sender log
cat data/sender_log.txt

# View receiver log
cat data/receiver_log.txt
```

**Log format:**
```
<timestamp> <event> <seqnum> <acknum> <flags> <payload_size>
```

**Event types:**
- `snd` - Segment sent
- `rcv` - Segment received
- `drp` - Dropped by PLC (Packet Loss & Corruption module)
- `crp` - Corrupted by PLC
- `dup` - Duplicate ACK received
- `rxt` - Retransmission (timeout)
- `frx` - Fast retransmit (3 duplicate ACKs)

## Docker Testing

Build and run in Docker:

```bash
# Build Docker image
docker-compose build

# Run test
docker-compose up

# Clean up
docker-compose down
```

## Troubleshooting

### Build Errors

**Error: `cannot find package`**
- Solution: Make sure you're in the repository root directory
- Run: `go mod tidy` to download dependencies

**Error: `module declares its path as urp-go-project`**
- Solution: The imports have been fixed to use `UDPAssignmentFun`
- Rebuild with: `go build ./cmd/sender` and `go build ./cmd/receiver`

### Test Failures

**Test fails with "timeout waiting for SYN-ACK"**
- Check firewall isn't blocking UDP ports 5000-6000
- Make sure no other process is using these ports

**High loss test (1.7, 1.8) fails frequently**
- This is expected - 50% loss can cause connection failures
- The protocol has retry limits to prevent infinite loops
- These tests verify graceful handling of extreme conditions

**All tests fail**
- Check executables built successfully: `ls -l sender.exe receiver.exe`
- Verify test files exist: `ls -l data/test_file.txt`
- Try building fresh: `go clean && go build ./cmd/...`

## Next Steps

- Read the [Chapter Guides](./chapters/) to understand how it works
- Examine the [source code](./internal/) to see the implementation
- Try the [examples](./examples/) in different programming languages
- Modify the protocol and run tests to see the effects!

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines on improving this educational resource.
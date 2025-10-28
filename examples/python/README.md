# Python UDP Examples

Simple UDP sender and receiver examples in Python.

## Requirements

- Python 3.9 or later (built-in socket library)
  - Python 3.8 and earlier are no longer supported (EOL)

## Running the Examples

### Terminal 1: Start the Receiver
```bash
cd examples/python
python receiver.py
```

Output:
```
Simple UDP Receiver listening on port 8080
Waiting for messages...
```

### Terminal 2: Run the Sender
```bash
cd examples/python
python sender.py
```

Output:
```
Sending message to server...
Received response: 'Echo: Hello, URP!'
```

## Key Features

- **Timeout handling** - Sender waits max 5 seconds for response
- **Echo server** - Receiver echoes back messages with "Echo: " prefix
- **Error handling** - Proper exception handling for network errors
- **Clean shutdown** - Receiver can be stopped with Ctrl+C

## Code Highlights

### UDP Socket Creation
```python
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
```

### Sending Data
```python
sock.sendto(b"Hello, URP!", ("127.0.0.1", 8080))
```

### Receiving Data
```python
data, addr = sock.recvfrom(1024)
message = data.decode()
```

## Learning Resources

- [Python socket documentation](https://docs.python.org/3/library/socket.html)
- [Socket Programming HOWTO](https://docs.python.org/3/howto/sockets.html)

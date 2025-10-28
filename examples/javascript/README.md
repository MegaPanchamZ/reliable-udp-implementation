# JavaScript (Node.js) UDP Examples

Simple UDP sender and receiver examples using Node.js dgram module.

## Requirements

- Node.js 18 or later (LTS version recommended)
  - Node.js 12-17 are no longer supported (EOL)

## Running the Examples

### Terminal 1: Start the Receiver
```bash
cd examples/javascript
node receiver.js
```

Output:
```
Simple UDP Receiver listening on 0.0.0.0:8080
Waiting for messages...
```

### Terminal 2: Run the Sender
```bash
cd examples/javascript
node sender.js
```

Output:
```
Sending message to server...
Received response: 'Echo: Hello, URP!'
```

## Key Features

- **Event-driven** - Uses Node.js event emitters
- **Timeout handling** - Sender has 5-second timeout
- **Echo server** - Receiver echoes with "Echo: " prefix
- **Error handling** - Proper error event handlers

## Code Highlights

### Creating UDP Socket
```javascript
const dgram = require('dgram');
const socket = dgram.createSocket('udp4');
```

### Sending Data
```javascript
socket.send(message, port, 'localhost', callback);
```

### Receiving Data
```javascript
socket.on('message', (msg, rinfo) => {
  console.log(`Received: '${msg}' from ${rinfo.address}:${rinfo.port}`);
});
```

## Learning Resources

- [Node.js dgram documentation](https://nodejs.org/api/dgram.html)
- [UDP Networking in Node.js](https://nodejs.org/en/knowledge/advanced/udp-multicast/)

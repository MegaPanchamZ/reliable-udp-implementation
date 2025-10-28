// --- Receiver ---
const dgram = require('dgram');
const server = dgram.createSocket('udp4');

server.on('message', (msg, rinfo) => {
  console.log(`Received: '${msg}' from ${rinfo.address}:${rinfo.port}`);
  const response = `Echo: ${msg}`;
  server.send(response, rinfo.port, rinfo.address, (err) => {
    if (err) console.error('Error sending response:', err);
  });
});

server.on('listening', () => {
  const addr = server.address();
  console.log(`Simple UDP Receiver listening on ${addr.address}:${addr.port}`);
  console.log('Waiting for messages...');
});

server.on('error', (err) => {
  console.error('Server error:', err);
  server.close();
});

server.bind(8080);

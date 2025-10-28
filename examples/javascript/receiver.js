// --- Receiver ---
const dgram = require('dgram');
const server = dgram.createSocket('udp4');
server.on('message', (msg, rinfo) => {
  console.log(`Received '${msg}' from ${rinfo.address}:${rinfo.port}`);
  server.send('Message received', rinfo.port, rinfo.address);
});
server.on('listening', () => console.log(`Listening on ${server.address().port}...`));
server.bind(8080);

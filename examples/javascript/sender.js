// --- Sender ---
const dgram = require('dgram');
const client = dgram.createSocket('udp4');

console.log('Sending message to server...');

const message = Buffer.from('Hello, URP!');
client.send(message, 8080, 'localhost', (err) => {
  if (err) {
    console.error('Error sending message:', err);
    client.close();
  }
});

client.on('message', (msg, rinfo) => {
  console.log(`Received response: '${msg}'`);
  client.close();
});

client.on('error', (err) => {
  console.error('Socket error:', err);
  client.close();
});

// Timeout after 5 seconds
setTimeout(() => {
  console.log('Timeout: No response received');
  client.close();
}, 5000);

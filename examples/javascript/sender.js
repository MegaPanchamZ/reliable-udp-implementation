// --- Sender ---
const dgram = require('dgram');
const client = dgram.createSocket('udp4');
client.send(Buffer.from('Hello, URP!'), 8080, 'localhost', (err) => {
  if (err) client.close();
});
client.on('message', (msg, rinfo) => {
  console.log(`Received response: '${msg}'`);
  client.close();
});

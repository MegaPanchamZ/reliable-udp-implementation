// --- Receiver ---
import 'dart:io';
import 'dart:convert';

void main() async {
  final socket = await RawDatagramSocket.bind(InternetAddress.anyIPv4, 8080);
  print('Listening on ${socket.address.address}:${socket.port}...');

  await for (final event in socket) {
    if (event == RawSocketEvent.read) {
      final datagram = socket.receive();
      if (datagram == null) continue;
      final message = utf8.decode(datagram.data);
      print("Received '$message' from ${datagram.address.address}:${datagram.port}");
      socket.send(utf8.encode('Message received'), datagram.address, datagram.port);
    }
  }
}

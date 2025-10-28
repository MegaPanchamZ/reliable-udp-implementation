// --- Sender ---
import 'dart:io';
import 'dart:convert';

void main() async {
  final destinationAddress = InternetAddress('127.0.0.1');
  final socket = await RawDatagramSocket.bind(InternetAddress.anyIPv4, 0);

  print('Sending message to server...');
  socket.send(utf8.encode('Hello, URP!'), destinationAddress, 8080);

  await for (final event in socket) {
    if (event == RawSocketEvent.read) {
      final datagram = socket.receive();
      if (datagram == null) continue;
      print("Received response: '${utf8.decode(datagram.data)}'");
      socket.close();
      break;
    }
  }
}

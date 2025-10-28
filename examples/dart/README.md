# Dart UDP Examples

Simple UDP sender and receiver examples in Dart using dart:io.

## Requirements

- Dart SDK 3.0 or later
  - Dart 2.x is no longer recommended (use 3.x for null safety and modern features)

## Installing Dart

### Windows (Chocolatey)
```bash
choco install dart-sdk
```

### macOS (Homebrew)
```bash
brew tap dart-lang/dart
brew install dart
```

### Linux
```bash
sudo apt-get update
sudo apt-get install dart
```

Or download from [dart.dev](https://dart.dev/get-dart)

## Running the Examples

### Terminal 1: Start the Receiver
```bash
cd examples/dart
dart receiver.dart
```

Output:
```
Listening on 0.0.0.0:8080...
```

### Terminal 2: Run the Sender
```bash
cd examples/dart
dart sender.dart
```

Output:
```
Sending message to server...
Received response: 'Echo: Hello, URP!'
```

## Key Features

- **Async/await** - Modern asynchronous programming
- **Stream-based** - Event-driven socket handling
- **UTF-8 encoding** - Built-in string encoding
- **Cross-platform** - Runs on Windows, macOS, Linux

## Code Highlights

### Creating Socket
```dart
final socket = await RawDatagramSocket.bind(InternetAddress.anyIPv4, 0);
```

### Sending Data
```dart
socket.send(utf8.encode('Hello, URP!'), destinationAddress, 8080);
```

### Receiving Data (Event Stream)
```dart
await for (final event in socket) {
  if (event == RawSocketEvent.read) {
    final datagram = socket.receive();
    if (datagram != null) {
      final message = utf8.decode(datagram.data);
      print("Received: '$message'");
    }
  }
}
```

## Code Structure

Dart uses event-driven programming with streams:
- `RawSocketEvent.read` - Data available to read
- `RawSocketEvent.write` - Socket ready to write
- `RawSocketEvent.closed` - Socket closed

## Learning Resources

- [Dart dart:io documentation](https://api.dart.dev/stable/dart-io/dart-io-library.html)
- [RawDatagramSocket API](https://api.dart.dev/stable/dart-io/RawDatagramSocket-class.html)
- [Dart Language Tour](https://dart.dev/guides/language/language-tour)

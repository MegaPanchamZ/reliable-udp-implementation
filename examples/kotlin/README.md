# Kotlin UDP Examples

Simple UDP sender and receiver examples in Kotlin using java.net.

## Requirements

- Kotlin compiler (kotlinc)
- JDK 8 or later

## Installing Kotlin

### SDKMAN (Linux/macOS)
```bash
sdk install kotlin
```

### Homebrew (macOS)
```bash
brew install kotlin
```

### Download
[kotlinlang.org/docs/command-line.html](https://kotlinlang.org/docs/command-line.html)

## Compiling

```bash
# Compile sender
kotlinc Sender.kt -include-runtime -d sender.jar

# Compile receiver
kotlinc Receiver.kt -include-runtime -d receiver.jar
```

## Running the Examples

### Terminal 1: Start the Receiver
```bash
cd examples/kotlin
java -jar receiver.jar
```

Output:
```
Listening on port 8080...
```

### Terminal 2: Run the Sender
```bash
cd examples/kotlin
java -jar sender.jar
```

Output:
```
Sending message to server...
Received response: 'Echo: Hello, URP!'
```

## Quick Run (No Compilation)

```bash
# Run directly with kotlin command
kotlin Sender.kt
kotlin Receiver.kt
```

## Key Features

- **JVM interop** - Uses Java networking APIs
- **Null safety** - Kotlin's type system prevents null errors
- **Concise syntax** - Less boilerplate than Java
- **Extension functions** - Enhanced ByteArray and String APIs

## Code Highlights

### Creating Socket
```kotlin
val socket = DatagramSocket()
```

### Sending Data
```kotlin
val message = "Hello, URP!".toByteArray()
val packet = DatagramPacket(message, message.size, serverAddress, 8080)
socket.send(packet)
```

### Receiving Data
```kotlin
val buffer = ByteArray(1024)
val packet = DatagramPacket(buffer, buffer.size)
socket.receive(packet)
val message = String(packet.data, 0, packet.length)
```

## Learning Resources

- [Kotlin Documentation](https://kotlinlang.org/docs/home.html)
- [Java DatagramSocket](https://docs.oracle.com/javase/8/docs/api/java/net/DatagramSocket.html)

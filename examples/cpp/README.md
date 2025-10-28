# C++ UDP Examples

Simple UDP sender and receiver examples in C++ using Boost.Asio.

## Requirements

- C++11 or later
- Boost library (Asio)

## Installing Boost

### Ubuntu/Debian
```bash
sudo apt-get install libboost-all-dev
```

### macOS (Homebrew)
```bash
brew install boost
```

### Windows
Download from [boost.org](https://www.boost.org/)

## Compiling

```bash
# Compile sender
g++ sender.cpp -o sender -lboost_system -std=c++11

# Compile receiver  
g++ receiver.cpp -o receiver -lboost_system -std=c++11
```

## Running the Examples

### Terminal 1: Start the Receiver
```bash
cd examples/cpp
./receiver
```

Output:
```
Listening on port 8080...
```

### Terminal 2: Run the Sender
```bash
cd examples/cpp
./sender
```

Output:
```
Sending message to server...
Received response: 'Echo: Hello, URP!'
```

## Key Features

- **Boost.Asio** - Modern asynchronous I/O library
- **Exception handling** - C++ try-catch for errors
- **RAII** - Automatic resource management
- **Type safety** - Strong type system

## Code Highlights

### Creating Socket
```cpp
boost::asio::io_context io_context;
udp::socket socket(io_context);
socket.open(udp::v4());
```

### Resolving Address
```cpp
udp::resolver resolver(io_context);
udp::endpoint receiver_endpoint = 
    *resolver.resolve(udp::v4(), "localhost", "8080").begin();
```

### Sending Data
```cpp
socket.send_to(boost::asio::buffer("Hello, URP!"), receiver_endpoint);
```

### Receiving Data
```cpp
char reply[1024];
udp::endpoint sender_endpoint;
size_t length = socket.receive_from(
    boost::asio::buffer(reply, 1024), sender_endpoint);
```

## Learning Resources

- [Boost.Asio Documentation](https://www.boost.org/doc/libs/release/doc/html/boost_asio.html)
- [Boost.Asio Tutorial](https://www.boost.org/doc/libs/release/doc/html/boost_asio/tutorial.html)

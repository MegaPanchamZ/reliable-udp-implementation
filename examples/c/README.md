# C UDP Examples

Simple UDP sender and receiver examples in C using BSD sockets (POSIX) and Winsock2 (Windows).

## Requirements

- **Linux/macOS:** GCC or Clang compiler
- **Windows:** MinGW-w64, MSVC, or Cygwin

## Compiling

### Linux/macOS
```bash
# Compile sender
gcc sender.c -o sender

# Compile receiver
gcc receiver.c -o receiver
```

### Windows (MinGW)
```bash
# Compile sender
gcc sender.c -o sender.exe -lws2_32

# Compile receiver
gcc receiver.c -o receiver.exe -lws2_32
```

### Windows (MSVC - Visual Studio Developer Command Prompt)
```bash
# Compile sender
cl sender.c ws2_32.lib

# Compile receiver
cl receiver.c ws2_32.lib
```

## Running the Examples

### Terminal 1: Start the Receiver
```bash
cd examples/c
./receiver          # Linux/macOS
# or
receiver.exe        # Windows
```

Output:
```
Listening on port 8080...
Press Ctrl+C to stop...
```

### Terminal 2: Run the Sender
```bash
cd examples/c
./sender            # Linux/macOS
# or
sender.exe          # Windows
```

Output:
```
Sending message to server...
Received response: 'Echo: Hello, URP!'
```

## Key Features

- **Cross-platform** - Works on Linux, macOS, and Windows
- **Low-level sockets** - Direct BSD socket API (POSIX) and Winsock2 (Windows)
- **System calls** - socket(), bind(), sendto(), recvfrom()
- **Error handling** - Comprehensive error checking with perror()
- **Memory safety** - Proper buffer management and bounds checking

## Code Highlights

### Platform Detection
```c
#ifdef _WIN32
    // Windows-specific code (Winsock2)
    #include <winsock2.h>
    WSAStartup(MAKEWORD(2, 2), &wsaData);
#else
    // POSIX-specific code (BSD sockets)
    #include <sys/socket.h>
#endif
```

### Creating Socket
```c
int sockfd = socket(AF_INET, SOCK_DGRAM, 0);
if (sockfd < 0) {
    perror("socket creation failed");
    exit(EXIT_FAILURE);
}
```

### Binding (Receiver)
```c
struct sockaddr_in servaddr;
servaddr.sin_family = AF_INET;
servaddr.sin_addr.s_addr = INADDR_ANY;
servaddr.sin_port = htons(PORT);

if (bind(sockfd, (struct sockaddr *)&servaddr, sizeof(servaddr)) < 0) {
    perror("bind failed");
    exit(EXIT_FAILURE);
}
```

### Sending Data (with error checking)
```c
if (sendto(sockfd, message, strlen(message), 0,
           (struct sockaddr *)&servaddr, sizeof(servaddr)) < 0) {
    perror("sendto failed");
    exit(EXIT_FAILURE);
}
```

### Receiving Data (with error checking)
```c
socklen_t len = sizeof(cliaddr);
int n = recvfrom(sockfd, buffer, MAX_BUFFER - 1, 0,
                 (struct sockaddr *)&cliaddr, &len);
if (n < 0) {
    perror("recvfrom failed");
    exit(EXIT_FAILURE);
}
buffer[n] = '\0';  // Null-terminate the received data
```

## Platform-Specific Notes

### Windows
- Must call `WSAStartup()` before using sockets
- Must call `WSACleanup()` before program exits
- Link with `ws2_32.lib` or `-lws2_32`
- Use `closesocket()` instead of `close()`

### Linux/macOS
- No initialization required
- Use standard POSIX socket functions
- Use `close()` to close sockets

## Learning Resources

- [Beej's Guide to Network Programming](https://beej.us/guide/bgnet/) - Excellent cross-platform guide
- [BSD Socket API (Linux)](https://man7.org/linux/man-pages/man7/socket.7.html)
- [Winsock2 Documentation (Windows)](https://docs.microsoft.com/en-us/windows/win32/winsock/)

// --- Sender ---
// Cross-platform UDP sender example
//
// Compile on Linux/macOS:
//   gcc sender.c -o sender
//
// Compile on Windows (MinGW):
//   gcc sender.c -o sender.exe -lws2_32
//
// Compile on Windows (MSVC):
//   cl sender.c ws2_32.lib

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Platform-specific headers
#ifdef _WIN32
    #include <winsock2.h>
    #include <ws2tcpip.h>
    #pragma comment(lib, "ws2_32.lib")

    // Windows compatibility definitions
    typedef int socklen_t;
    #define close closesocket
#else
    #include <sys/socket.h>
    #include <netinet/in.h>
    #include <arpa/inet.h>
    #include <unistd.h>
#endif

#define PORT 8080
#define MAX_BUFFER 1024

int main() {
#ifdef _WIN32
    // Initialize Winsock on Windows
    WSADATA wsaData;
    if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0) {
        fprintf(stderr, "WSAStartup failed\n");
        return EXIT_FAILURE;
    }
#endif

    int sockfd;
    char buffer[MAX_BUFFER];
    const char *message = "Hello, URP!";
    struct sockaddr_in servaddr;

    // Create socket
    if ((sockfd = socket(AF_INET, SOCK_DGRAM, 0)) < 0) {
        perror("socket creation failed");
#ifdef _WIN32
        WSACleanup();
#endif
        exit(EXIT_FAILURE);
    }

    // Setup server address
    memset(&servaddr, 0, sizeof(servaddr));
    servaddr.sin_family = AF_INET;
    servaddr.sin_port = htons(PORT);
    servaddr.sin_addr.s_addr = inet_addr("127.0.0.1");

    // Send message
    printf("Sending message to server...\n");
    if (sendto(sockfd, message, strlen(message), 0,
               (const struct sockaddr *)&servaddr, sizeof(servaddr)) < 0) {
        perror("sendto failed");
        close(sockfd);
#ifdef _WIN32
        WSACleanup();
#endif
        exit(EXIT_FAILURE);
    }

    // Receive response
    socklen_t len = sizeof(servaddr);
    int n = recvfrom(sockfd, buffer, MAX_BUFFER - 1, 0,
                     (struct sockaddr *)&servaddr, &len);
    if (n < 0) {
        perror("recvfrom failed");
        close(sockfd);
#ifdef _WIN32
        WSACleanup();
#endif
        exit(EXIT_FAILURE);
    }

    buffer[n] = '\0';
    printf("Received response: '%s'\n", buffer);

    // Cleanup
    close(sockfd);
#ifdef _WIN32
    WSACleanup();
#endif
    return 0;
}

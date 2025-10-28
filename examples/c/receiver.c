// --- Receiver ---
// Cross-platform UDP receiver example
//
// Compile on Linux/macOS:
//   gcc receiver.c -o receiver
//
// Compile on Windows (MinGW):
//   gcc receiver.c -o receiver.exe -lws2_32
//
// Compile on Windows (MSVC):
//   cl receiver.c ws2_32.lib

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
    struct sockaddr_in servaddr, cliaddr;

    // Create socket file descriptor
    if ((sockfd = socket(AF_INET, SOCK_DGRAM, 0)) < 0) {
        perror("socket creation failed");
#ifdef _WIN32
        WSACleanup();
#endif
        exit(EXIT_FAILURE);
    }

    // Setup server address
    memset(&servaddr, 0, sizeof(servaddr));
    memset(&cliaddr, 0, sizeof(cliaddr));

    servaddr.sin_family = AF_INET;
    servaddr.sin_addr.s_addr = INADDR_ANY;
    servaddr.sin_port = htons(PORT);

    // Bind socket to port
    if (bind(sockfd, (const struct sockaddr *)&servaddr, sizeof(servaddr)) < 0) {
        perror("bind failed");
        close(sockfd);
#ifdef _WIN32
        WSACleanup();
#endif
        exit(EXIT_FAILURE);
    }

    printf("Listening on port %d...\n", PORT);
    printf("Press Ctrl+C to stop...\n");

    socklen_t len = sizeof(cliaddr);
    while (1) {
        // Receive message
        int n = recvfrom(sockfd, buffer, MAX_BUFFER - 1, 0,
                        (struct sockaddr *)&cliaddr, &len);
        if (n < 0) {
            perror("recvfrom failed");
            continue;
        }

        buffer[n] = '\0';
        printf("Received '%s'\n", buffer);

        // Echo back with prefix
        char response[MAX_BUFFER];
        snprintf(response, MAX_BUFFER, "Echo: %s", buffer);

        if (sendto(sockfd, response, strlen(response), 0,
                   (const struct sockaddr *)&cliaddr, len) < 0) {
            perror("sendto failed");
        }
    }

    // Cleanup (unreachable in this example, but good practice)
    close(sockfd);
#ifdef _WIN32
    WSACleanup();
#endif
    return 0;
}

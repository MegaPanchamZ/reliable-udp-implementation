// --- Sender ---
// (Compile with: gcc sender.c -o sender)
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <unistd.h>

#define PORT 8080
#define MAX_BUFFER 1024

int main() {
    int sockfd;
    char buffer[MAX_BUFFER];
    const char *message = "Hello, URP!";
    struct sockaddr_in servaddr;

    if ((sockfd = socket(AF_INET, SOCK_DGRAM, 0)) < 0) {
        perror("socket creation failed");
        exit(EXIT_FAILURE);
    }

    memset(&servaddr, 0, sizeof(servaddr));
    servaddr.sin_family = AF_INET;
    servaddr.sin_port = htons(PORT);
    servaddr.sin_addr.s_addr = INADDR_ANY; // Connect to localhost

    sendto(sockfd, message, strlen(message), 0, (const struct sockaddr *)&servaddr, sizeof(servaddr));
    printf("Sending message to server...\n");

    socklen_t len = sizeof(servaddr);
    int n = recvfrom(sockfd, (char *)buffer, MAX_BUFFER, 0, (struct sockaddr *)&servaddr, &len);
    buffer[n] = '\0';
    printf("Received response: '%s'\n", buffer);

    close(sockfd);
    return 0;
}

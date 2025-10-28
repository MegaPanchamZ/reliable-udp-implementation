// --- Receiver ---
// (Compile with: gcc receiver.c -o receiver)
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
    struct sockaddr_in servaddr, cliaddr;

    // Create socket file descriptor
    if ((sockfd = socket(AF_INET, SOCK_DGRAM, 0)) < 0) {
        perror("socket creation failed");
        exit(EXIT_FAILURE);
    }

    memset(&servaddr, 0, sizeof(servaddr));
    memset(&cliaddr, 0, sizeof(cliaddr));

    servaddr.sin_family = AF_INET;
    servaddr.sin_addr.s_addr = INADDR_ANY;
    servaddr.sin_port = htons(PORT);

    if (bind(sockfd, (const struct sockaddr *)&servaddr, sizeof(servaddr)) < 0) {
        perror("bind failed");
        exit(EXIT_FAILURE);
    }

    printf("Listening on port %d...\n", PORT);
    socklen_t len = sizeof(cliaddr);
    while (1) {
        int n = recvfrom(sockfd, (char *)buffer, MAX_BUFFER, 0, (struct sockaddr *)&cliaddr, &len);
        buffer[n] = '\0';
        printf("Received '%s'\n", buffer);
        sendto(sockfd, "Message received", strlen("Message received"), 0, (const struct sockaddr *)&cliaddr, len);
    }
    close(sockfd);
    return 0;
}

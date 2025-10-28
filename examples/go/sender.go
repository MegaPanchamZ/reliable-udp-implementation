// --- Sender ---
package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error creating connection:", err)
		return
	}
	defer conn.Close()

	// Set read timeout
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	message := "Hello, URP!"
	fmt.Println("Sending message to server...")
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		return
	}

	fmt.Printf("Received response: '%s'\n", string(buffer[:n]))
}

// --- Receiver ---
package main

import (
	"fmt"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	fmt.Println("Simple UDP Receiver listening on :8080")
	fmt.Println("Waiting for messages...")

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("Received: '%s' from %s\n", message, remoteAddr)

		// Echo back with prefix
		response := "Echo: " + message
		_, err = conn.WriteToUDP([]byte(response), remoteAddr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}
	}
}

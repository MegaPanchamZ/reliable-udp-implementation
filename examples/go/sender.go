// --- Sender ---
package main
import ("fmt"; "net")
func main() {
	serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	conn, _ := net.DialUDP("udp", nil, serverAddr)
	defer conn.Close()
	fmt.Println("Sending message to server...")
	conn.Write([]byte("Hello, URP!"))
	buffer := make([]byte, 1024)
	n, _, _ := conn.Read(buffer)
	fmt.Printf("Received response: '%s'\n", string(buffer[:n]))
}


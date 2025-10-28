// --- Receiver ---
package main
import ("fmt"; "net")
func main() {
	addr, _ := net.ResolveUDPAddr("udp", ":8080")
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()
	buffer := make([]byte, 1024)
	fmt.Println("Listening on port 8080...")
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil { continue }
		fmt.Printf("Received '%s' from %s\n", string(buffer[:n]), remoteAddr)
		conn.WriteToUDP([]byte("Message received"), remoteAddr)
	}
}


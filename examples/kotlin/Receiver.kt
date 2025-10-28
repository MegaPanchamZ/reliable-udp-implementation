// --- Receiver ---
// (Compile with kotlinc Receiver.kt -include-runtime -d receiver.jar)
import java.net.DatagramPacket
import java.net.DatagramSocket

fun main() {
    val socket = DatagramSocket(8080)
    val buffer = ByteArray(1024)
    println("Listening on port 8080...")

    while (true) {
        val packet = DatagramPacket(buffer, buffer.size)
        socket.receive(packet)
        val message = String(packet.data, 0, packet.length)
        println("Received '$message' from ${packet.address}:${packet.port}")

        // Echo back with prefix
        val response = "Echo: $message".toByteArray()
        val responsePacket = DatagramPacket(response, response.size, packet.address, packet.port)
        socket.send(responsePacket)
    }
}

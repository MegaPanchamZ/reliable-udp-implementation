// --- Sender ---
// (Compile with kotlinc Sender.kt -include-runtime -d sender.jar)
import java.net.DatagramPacket
import java.net.DatagramSocket
import java.net.InetAddress

fun main() {
    val socket = DatagramSocket()
    val serverAddress = InetAddress.getByName("localhost")
    val message = "Hello, URP!".toByteArray()

    println("Sending message to server...")
    val packet = DatagramPacket(message, message.size, serverAddress, 8080)
    socket.send(packet)

    val buffer = ByteArray(1024)
    val responsePacket = DatagramPacket(buffer, buffer.size)
    socket.receive(responsePacket)
    val response = String(responsePacket.data, 0, responsePacket.length)
    println("Received response: '$response'")
    socket.close()
}

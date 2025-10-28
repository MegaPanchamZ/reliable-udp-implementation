# --- Sender ---
import socket
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.sendto(b"Hello, URP!", ("127.0.0.1", 8080))
data, server = sock.recvfrom(1024)
print(f"Received response: '{data.decode()}'")
sock.close()

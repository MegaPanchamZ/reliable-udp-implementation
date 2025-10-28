# --- Sender ---
import socket

print("Sending message to server...")
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.settimeout(5.0)  # 5 second timeout

try:
    sock.sendto(b"Hello, URP!", ("127.0.0.1", 8080))
    data, server = sock.recvfrom(1024)
    print(f"Received response: '{data.decode()}'")
except socket.timeout:
    print("Error: No response from server (timeout)")
except Exception as e:
    print(f"Error: {e}")
finally:
    sock.close()

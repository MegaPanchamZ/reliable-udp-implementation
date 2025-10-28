# --- Receiver ---
import socket

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.bind(("0.0.0.0", 8080))
print("Simple UDP Receiver listening on port 8080")
print("Waiting for messages...")

try:
    while True:
        data, addr = sock.recvfrom(1024)
        message = data.decode()
        print(f"Received: '{message}' from {addr}")
        response = f"Echo: {message}"
        sock.sendto(response.encode(), addr)
except KeyboardInterrupt:
    print("\nShutting down...")
finally:
    sock.close()

// --- Sender ---
// (Link with: g++ sender.cpp -o sender -lboost_system)
#include <iostream>
#include <boost/asio.hpp>

using boost::asio::ip::udp;

int main() {
    try {
        boost::asio::io_context io_context;
        udp::resolver resolver(io_context);
        udp::endpoint receiver_endpoint = *resolver.resolve(udp::v4(), "localhost", "8080").begin();
        udp::socket socket(io_context);
        socket.open(udp::v4());

        std::cout << "Sending message to server..." << std::endl;
        socket.send_to(boost::asio::buffer("Hello, URP!"), receiver_endpoint);

        char reply[1024];
        udp::endpoint sender_endpoint;
        size_t reply_length = socket.receive_from(boost::asio::buffer(reply, 1024), sender_endpoint);
        std::cout << "Received response: '";
        std::cout.write(reply, reply_length);
        std::cout << "'" << std::endl;
    } catch (std::exception& e) {
        std::cerr << e.what() << std::endl;
    }
    return 0;
}

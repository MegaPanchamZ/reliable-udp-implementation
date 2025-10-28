// --- Receiver ---
// (Requires Asio library. Link with: g++ receiver.cpp -o receiver -lboost_system)
#include <iostream>
#include <string>
#include <boost/asio.hpp>

using boost::asio::ip::udp;

int main() {
    try {
        boost::asio::io_context io_context;
        udp::socket socket(io_context, udp::endpoint(udp::v4(), 8080));
        std::cout << "Listening on port 8080..." << std::endl;

        while (true) {
            char data[1024];
            udp::endpoint sender_endpoint;
            size_t length = socket.receive_from(boost::asio::buffer(data, 1024), sender_endpoint);
            std::cout << "Received '" << std::string(data, length) << "'" << std::endl;
            socket.send_to(boost::asio::buffer("Message received"), sender_endpoint);
        }
    } catch (std::exception& e) {
        std::cerr << e.what() << std::endl;
    }
    return 0;
}

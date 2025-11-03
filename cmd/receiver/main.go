package main

import (
	"urp-go/internal/logger"
	"urp-go/internal/plc"
	"urp-go/internal/protocol"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Parse command-line flags
	receiverPort := flag.Int("receiver_port", 0, "Receiver's UDP port")
	senderPort := flag.Int("sender_port", 0, "Sender's UDP port (for reference)")
	filePath := flag.String("file", "", "Path to write the received text file")
	maxWin := flag.Int("max_win", 0, "Receive window size")
	rlp := flag.Float64("rlp", 0.0, "Reverse loss probability")
	rcp := flag.Float64("rcp", 0.0, "Reverse corruption probability")
	logPath := flag.String("log", "data/receiver_log.txt", "Path to the log file")

	flag.Parse()

	// Validate required flags
	if *receiverPort == 0 || *senderPort == 0 || *filePath == "" || *maxWin == 0 {
		fmt.Fprintf(os.Stderr, "Usage: receiver -receiver_port=<port> -sender_port=<port> -file=<path> -max_win=<bytes>\n")
		os.Exit(1)
	}

	// Create logger
	log, err := logger.New(*logPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	// Create PLC module (only reverse path for receiver)
	plcModule := plc.New(0.0, *rlp, 0.0, *rcp, log)

	// Create receiver
	receiver, err := protocol.NewReceiver(*receiverPort, *senderPort, *filePath, *maxWin, log, plcModule)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating receiver: %v\n", err)
		os.Exit(1)
	}
	defer receiver.Close()

	// Start listening
	fmt.Printf("Receiver listening on port %d\n", *receiverPort)
	fmt.Printf("Output file: %s\n", *filePath)
	fmt.Printf("Max window: %d bytes\n", *maxWin)

	if err := receiver.Listen(); err != nil {
		fmt.Fprintf(os.Stderr, "Error during reception: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File received successfully!")
	log.Flush()
}

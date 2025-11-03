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
	senderPort := flag.Int("sender_port", 0, "Sender's UDP port")
	receiverHost := flag.String("receiver_host", "127.0.0.1", "Receiver's hostname or IP")
	receiverPort := flag.Int("receiver_port", 0, "Receiver's UDP port")
	filePath := flag.String("file", "", "Path to the text file to send")
	maxWin := flag.Int("max_win", 0, "Maximum window size in bytes")
	rto := flag.Int("rto", 0, "Retransmission timeout in milliseconds")
	flp := flag.Float64("flp", 0.0, "Forward loss probability")
	rlp := flag.Float64("rlp", 0.0, "Reverse loss probability")
	fcp := flag.Float64("fcp", 0.0, "Forward corruption probability")
	rcp := flag.Float64("rcp", 0.0, "Reverse corruption probability")
	logPath := flag.String("log", "data/sender_log.txt", "Path to the log file")

	flag.Parse()

	// Validate required flags
	if *senderPort == 0 || *receiverPort == 0 || *filePath == "" || *maxWin == 0 || *rto == 0 {
		fmt.Fprintf(os.Stderr, "Usage: sender -sender_port=<port> -receiver_host=<host> -receiver_port=<port> -file=<path> -max_win=<bytes> -rto=<ms> -flp=<prob> -rlp=<prob> -fcp=<prob> -rcp=<prob>\n")
		os.Exit(1)
	}

	// Read the input file
	fileData, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Create logger
	log, err := logger.New(*logPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	// Create PLC module
	plcModule := plc.New(*flp, *rlp, *fcp, *rcp, log)

	// Create sender
	sender, err := protocol.NewSender(*senderPort, *receiverHost, *receiverPort, *maxWin, *rto, log, plcModule)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating sender: %v\n", err)
		os.Exit(1)
	}
	defer sender.Close()

	// Send the file
	fmt.Printf("Sending file: %s (%d bytes)\n", *filePath, len(fileData))
	fmt.Printf("Sender port: %d, Receiver: %s:%d\n", *senderPort, *receiverHost, *receiverPort)
	fmt.Printf("Max window: %d, RTO: %dms\n", *maxWin, *rto)
	fmt.Printf("Loss probabilities - Forward: %.2f, Reverse: %.2f\n", *flp, *rlp)
	fmt.Printf("Corruption probabilities - Forward: %.2f, Reverse: %.2f\n", *fcp, *rcp)

	if err := sender.SendFile(fileData); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File sent successfully!")
	log.Flush()
}

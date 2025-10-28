#!/bin/bash

# Cross-platform test runner for the URP protocol implementation.
# Works on macOS, Linux, and Windows (with Git Bash or WSL).

# --- Configuration ---
SENDER_PORT=9090
RECEIVER_PORT=9091
TIMEOUT_SECONDS=30 # Max time to wait for a single test case
TEST_FILE="data/test_file.txt"
RECEIVED_FILE="data/received_file.txt"
LOG_DIR="data/logs"
SENDER_LOG="$LOG_DIR/sender.log"
RECEIVER_LOG="$LOG_DIR/receiver.log"

# --- Colors for output ---
COLOR_GREEN='\033[0;32m'
COLOR_RED='\033[0;31m'
COLOR_YELLOW='\033[0;33m'
COLOR_NC='\033[0m' # No Color

# --- Helper Functions ---
print_test_header() {
    echo -e "\n${COLOR_YELLOW}--------------------------------------------------${COLOR_NC}"
    echo -e "${COLOR_YELLOW}Running Test: $1${COLOR_NC}"
    echo -e "${COLOR_YELLOW}Description: $2${COLOR_NC}"
    echo -e "${COLOR_YELLOW}--------------------------------------------------${COLOR_NC}"
}

check_diff() {
    if diff -q "$TEST_FILE" "$RECEIVED_FILE" >/dev/null; then
        echo -e "${COLOR_GREEN}PASS: Received file is identical to the original.${COLOR_NC}"
        return 0
    else
        echo -e "${COLOR_RED}FAIL: Received file differs from the original.${COLOR_NC}"
        return 1
    fi
}

cleanup() {
    # Kill background processes (receiver)
    # Using pkill is more robust than relying on a single PID
    pkill -f "./receiver"
    rm -f "$RECEIVED_FILE"
}

# --- Main Script ---

# 1. Initial Setup
echo "Setting up test environment..."
mkdir -p "$LOG_DIR"
rm -f "$LOG_DIR"/* # Clean old logs
rm -f "$RECEIVED_FILE"

# Create a test file
echo "Creating test file with random data..."
head -c 100K /dev/urandom | base64 > "$TEST_FILE"

# 2. Build the Go binaries
echo "Building sender and receiver..."
go build -o receiver ./src/receiver.go
go build -o sender ./src/sender.go

if [ $? -ne 0 ]; then
    echo -e "${COLOR_RED}Build failed. Aborting tests.${COLOR_NC}"
    exit 1
fi
echo "Build successful."

# Trap Ctrl+C and other signals to ensure cleanup
trap cleanup EXIT INT TERM

# --- Test Cases ---

# Test 1: Perfect World
print_test_header "Perfect World" "0% loss, 0% corruption. Verifies basic correctness."
./receiver -receiver_port $RECEIVER_PORT -file "$RECEIVED_FILE" -log "$RECEIVER_LOG" &
RECEIVER_PID=$!
sleep 1 # Give receiver time to start
./sender -sender_port $SENDER_PORT -receiver_port $RECEIVER_PORT -file "$TEST_FILE" -rto 500 -max_win 5000 -log "$SENDER_LOG"
wait $RECEIVER_PID
check_diff
cleanup

# Test 2: High Forward Loss
print_test_header "High Forward Loss" "20% forward loss. Tests sender's retransmission timer."
./receiver -receiver_port $RECEIVER_PORT -file "$RECEIVED_FILE" -log "$RECEIVER_LOG" &
RECEIVER_PID=$!
sleep 1
./sender -sender_port $SENDER_PORT -receiver_port $RECEIVER_PORT -file "$TEST_FILE" -rto 500 -max_win 5000 -flp 0.2 -log "$SENDER_LOG"
wait $RECEIVER_PID
check_diff
cleanup

# Test 3: High Reverse (ACK) Loss
print_test_header "High Reverse (ACK) Loss" "20% reverse loss. Also tests retransmission."
./receiver -receiver_port $RECEIVER_PORT -file "$RECEIVED_FILE" -rlp 0.2 -log "$RECEIVER_LOG" &
RECEIVER_PID=$!
sleep 1
./sender -sender_port $SENDER_PORT -receiver_port $RECEIVER_PORT -file "$TEST_FILE" -rto 500 -max_win 5000 -log "$SENDER_LOG"
wait $RECEIVER_PID
check_diff
cleanup

# Test 4: High Corruption
print_test_header "High Corruption" "15% forward corruption. Tests receiver's checksum validation."
./receiver -receiver_port $RECEIVER_PORT -file "$RECEIVED_FILE" -log "$RECEIVER_LOG" &
RECEIVER_PID=$!
sleep 1
./sender -sender_port $SENDER_PORT -receiver_port $RECEIVER_PORT -file "$TEST_FILE" -rto 500 -max_win 5000 -fcp 0.15 -log "$SENDER_LOG"
wait $RECEIVER_PID
check_diff
cleanup

# Test 5: The Chaos Test
print_test_header "Chaos Test" "10% forward loss, 10% reverse loss, 5% forward corruption."
./receiver -receiver_port $RECEIVER_PORT -file "$RECEIVED_FILE" -rlp 0.1 -log "$RECEIVER_LOG" &
RECEIVER_PID=$!
sleep 1
./sender -sender_port $SENDER_PORT -receiver_port $RECEIVER_PORT -file "$TEST_FILE" -rto 500 -max_win 5000 -flp 0.1 -fcp 0.05 -log "$SENDER_LOG"
wait $RECEIVER_PID
check_diff
cleanup

echo -e "\n${COLOR_GREEN}All tests completed.${COLOR_NC}"


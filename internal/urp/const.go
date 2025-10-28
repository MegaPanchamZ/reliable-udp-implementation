package urp

// Protocol constants

// MSS is the Maximum Segment Size (payload data only)
const MSS = 1000

// Header size in bytes
const HeaderSize = 6

// Flag bits
const (
	FlagACK = 1 << 0 // 0x01
	FlagSYN = 1 << 1 // 0x02
	FlagFIN = 1 << 2 // 0x04
)

// Sender states
const (
	StateCLOSED = iota
	StateSYNSENT
	StateESTABLISHED
	StateFINSENT
	StateFINACKED
)

// Receiver states
const (
	StateLISTEN = iota + 10 // Offset to avoid collision with sender states
	StateSYNRCVD
	// StateESTABLISHED is shared with sender (value 2)
	StateTIMEWAIT
)

// Timeout values
const (
	TimeWaitDuration = 2000 // 2 seconds in milliseconds
)

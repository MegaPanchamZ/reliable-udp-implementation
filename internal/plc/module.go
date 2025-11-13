package plc

import (
	"math/rand"
	"urp-go/internal/logger"
	"urp-go/internal/urp"
)

// Module implements Packet Loss and Corruption simulation
type Module struct {
	forwardLossProb       float64
	reverseLossProb       float64
	forwardCorruptionProb float64
	reverseCorruptionProb float64
	logger                *logger.Logger
	rng                   *rand.Rand
}

// New creates a new PLC module with the specified probabilities
func New(flp, rlp, fcp, rcp float64, log *logger.Logger) *Module {
	return &Module{
		forwardLossProb:       flp,
		reverseLossProb:       rlp,
		forwardCorruptionProb: fcp,
		reverseCorruptionProb: rcp,
		logger:                log,
		rng:                   rand.New(rand.NewSource(rand.Int63())),
	}
}

// ShouldDropForward determines if a forward packet should be dropped
func (m *Module) ShouldDropForward() bool {
	return m.rng.Float64() < m.forwardLossProb
}

// ShouldDropReverse determines if a reverse packet should be dropped
func (m *Module) ShouldDropReverse() bool {
	return m.rng.Float64() < m.reverseLossProb
}

// ShouldCorruptForward determines if a forward packet should be corrupted
func (m *Module) ShouldCorruptForward() bool {
	return m.rng.Float64() < m.forwardCorruptionProb
}

// ShouldCorruptReverse determines if a reverse packet should be corrupted
func (m *Module) ShouldCorruptReverse() bool {
	return m.rng.Float64() < m.reverseCorruptionProb
}

// ProcessOutgoingForward processes an outgoing forward packet (sender to receiver)
// Returns the packet data and whether it should be sent
func (m *Module) ProcessOutgoingForward(data []byte, seg *urp.URPSegment) ([]byte, bool) {
	// Check for drop first
	if m.ShouldDropForward() {
		if m.logger != nil {
			m.logger.LogDrop(seg.SeqNum)
		}
		return nil, false
	}

	// Check for corruption
	if m.ShouldCorruptForward() {
		corrupted := make([]byte, len(data))
		copy(corrupted, data)
		urp.CorruptData(corrupted)
		if m.logger != nil {
			m.logger.LogCorrupt(seg.SeqNum)
		}
		return corrupted, true
	}

	return data, true
}

// ProcessOutgoingReverse processes an outgoing reverse packet (receiver to sender, ACKs)
// Returns the packet data and whether it should be sent
func (m *Module) ProcessOutgoingReverse(data []byte, seg *urp.URPSegment) ([]byte, bool) {
	// For ACK segments, SeqNum contains the acknowledgment number
	ackNum := seg.SeqNum

	// Check for drop first
	if m.ShouldDropReverse() {
		if m.logger != nil {
			m.logger.LogDrop(ackNum)
		}
		return nil, false
	}

	// Check for corruption
	if m.ShouldCorruptReverse() {
		corrupted := make([]byte, len(data))
		copy(corrupted, data)
		urp.CorruptData(corrupted)
		if m.logger != nil {
			m.logger.LogCorrupt(ackNum)
		}
		return corrupted, true
	}

	return data, true
}

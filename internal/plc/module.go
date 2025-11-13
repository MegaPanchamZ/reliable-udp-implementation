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
// Returns the packet data, status, and whether it should be sent
func (m *Module) ProcessOutgoingForward(data []byte, seg *urp.URPSegment) ([]byte, logger.StatusType, bool) {
	// Check for drop first
	if m.ShouldDropForward() {
		if m.logger != nil {
			m.logger.IncrementStat("plc_forward_dropped", 1)
		}
		return nil, logger.StatusDrop, false
	}

	// Check for corruption
	if m.ShouldCorruptForward() {
		corrupted := make([]byte, len(data))
		copy(corrupted, data)
		urp.CorruptData(corrupted, m.rng)
		if m.logger != nil {
			m.logger.IncrementStat("plc_forward_corrupted", 1)
		}
		return corrupted, logger.StatusCor, true
	}

	return data, logger.StatusOK, true
}

// ProcessOutgoingReverse processes an outgoing reverse packet (receiver to sender, ACKs)
// Returns the packet data, status, and whether it should be sent
func (m *Module) ProcessOutgoingReverse(data []byte, seg *urp.URPSegment) ([]byte, logger.StatusType, bool) {
	// Check for drop first
	if m.ShouldDropReverse() {
		if m.logger != nil {
			m.logger.IncrementStat("plc_reverse_dropped", 1)
		}
		return nil, logger.StatusDrop, false
	}

	// Check for corruption
	if m.ShouldCorruptReverse() {
		corrupted := make([]byte, len(data))
		copy(corrupted, data)
		urp.CorruptData(corrupted, m.rng)
		if m.logger != nil {
			m.logger.IncrementStat("plc_reverse_corrupted", 1)
		}
		return corrupted, logger.StatusCor, true
	}

	return data, logger.StatusOK, true
}

// ProcessIncomingReverse processes an incoming reverse packet (ACKs from receiver)
// Returns the packet data, status, and whether it should be delivered
func (m *Module) ProcessIncomingReverse(data []byte) ([]byte, logger.StatusType, bool) {
	// Check for drop first
	if m.ShouldDropReverse() {
		if m.logger != nil {
			m.logger.IncrementStat("plc_reverse_dropped", 1)
		}
		return nil, logger.StatusDrop, false
	}

	// Check for corruption
	if m.ShouldCorruptReverse() {
		corrupted := make([]byte, len(data))
		copy(corrupted, data)
		urp.CorruptData(corrupted, m.rng)
		if m.logger != nil {
			m.logger.IncrementStat("plc_reverse_corrupted", 1)
		}
		return corrupted, logger.StatusCor, true
	}

	return data, logger.StatusOK, true
}

// ProcessIncomingForward processes an incoming forward packet (DATA/SYN/FIN from sender)
// Returns the packet data, status, and whether it should be delivered
func (m *Module) ProcessIncomingForward(data []byte) ([]byte, logger.StatusType, bool) {
	// Check for drop first
	if m.ShouldDropForward() {
		if m.logger != nil {
			m.logger.IncrementStat("plc_forward_dropped", 1)
		}
		return nil, logger.StatusDrop, false
	}

	// Check for corruption
	if m.ShouldCorruptForward() {
		corrupted := make([]byte, len(data))
		copy(corrupted, data)
		urp.CorruptData(corrupted, m.rng)
		if m.logger != nil {
			m.logger.IncrementStat("plc_forward_corrupted", 1)
		}
		return corrupted, logger.StatusCor, true
	}

	return data, logger.StatusOK, true
}

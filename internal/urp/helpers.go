package urp

import "urp-go/internal/logger"

// GetSegmentType returns the segment type for logging
func GetSegmentType(seg *URPSegment) logger.SegmentType {
	if seg.HasFlag(FlagACK) {
		return logger.SegmentACK
	}
	if seg.HasFlag(FlagSYN) {
		return logger.SegmentSYN
	}
	if seg.HasFlag(FlagFIN) {
		return logger.SegmentFIN
	}
	return logger.SegmentDATA
}


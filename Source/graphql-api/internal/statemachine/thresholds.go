package statemachine

import "fmt"

// Thresholds sets how many consecutive results are needed to change state.
type Thresholds struct {
	Failure  int // consecutive failures before going down
	Recovery int // consecutive successes before going up
}

// Validate reports whether both thresholds are at least 1.
func (t Thresholds) Validate() error {
	if t.Failure < 1 {
		return fmt.Errorf("invalid thresholds: %w", ErrThresholdsFailureTooLow)
	}
	if t.Recovery < 1 {
		return fmt.Errorf("invalid thresholds: %w", ErrThresholdsRecoveryTooLow)
	}
	return nil
}

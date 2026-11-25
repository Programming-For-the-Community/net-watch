package statemachine

import (
	"fmt"

	"github.com/Programming-For-the-Community/net-watch/graphql-api/internal/model"
)

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

// Status is the current state of a check and its current streak of
// consecutive failures or successes. At most one counter is non-zero.
type Status struct {
	State     model.CheckState
	Failures  int
	Successes int
}

// Validate reports whether s is a status that Next could have produced
// under thresholds t. Paused statuses are not compared against t, since
// thresholds may change while a check is paused.
func (s Status) Validate(t Thresholds) error {
	if s.State == "" {
		return fmt.Errorf("invalid status: %w", ErrStatusMissingStateField)
	}
	if s.Failures < 0 {
		return fmt.Errorf("invalid status: %w", ErrStatusFailuresFieldNegative)
	}
	if s.Successes < 0 {
		return fmt.Errorf("invalid status: %w", ErrStatusSuccessesFieldNegative)
	}
	if s.Failures > 0 && s.Successes > 0 {
		return fmt.Errorf("invalid status: %w", ErrInvalidFailuresAndSuccessesCounts)
	}

	switch s.State {
	case model.StatePaused:
		return nil
	case model.StateUnknown:
		if s.Failures >= t.Failure {
			return fmt.Errorf("invalid status: %w", ErrInvalidUnknownFailuresStatus)
		}
		if s.Successes >= t.Recovery {
			return fmt.Errorf("invalid status: %w", ErrInvalidUnknownSuccessesStatus)
		}
	case model.StateUp:
		if s.Failures >= t.Failure {
			return fmt.Errorf("invalid status: %w", ErrInvalidUpStatus)
		}
	case model.StateDown:
		if s.Successes >= t.Recovery {
			return fmt.Errorf("invalid status: %w", ErrInvalidDownStatus)
		}
	default:
		return fmt.Errorf("invalid status %q: %w", s.State, ErrStatusUnknownState)
	}
	return nil
}

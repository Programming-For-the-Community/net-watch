package statemachine

import (
	"github.com/Programming-For-the-Community/net-watch/graphql-api/internal/model"
)

// Next applies one check result to s and returns the new status. The state
// changes on the result that brings the streak up to the threshold: the
// t.Failure-th consecutive failure goes down, and the t.Recovery-th
// consecutive success goes up. A paused status is returned unchanged.
//
// If s or t is invalid, Next returns an empty Status and an error.
func Next(s Status, success bool, t Thresholds) (Status, error) {
	if err := t.Validate(); err != nil {
		return Status{}, err
	}
	if err := s.Validate(t); err != nil {
		return Status{}, err
	}

	if s.State == model.StatePaused {
		return s, nil
	}

	next := s
	if success {
		next.Failures = 0
		next.Successes++
		if next.Successes >= t.Recovery {
			next.State = model.StateUp
		}
	} else {
		next.Successes = 0
		next.Failures++
		if next.Failures >= t.Failure {
			next.State = model.StateDown
		}
	}

	return next, nil
}

// Reset returns a status with no history: unknown, with both counters at 0.
// Use it when a check is created, resumed after a pause, or reconfigured.
func Reset() Status {
	return Status{State: model.StateUnknown}
}

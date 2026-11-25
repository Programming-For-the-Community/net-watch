package statemachine

import (
	"fmt"

	"github.com/Programming-For-the-Community/net-watch/graphql-api/internal/model"
)

func Next(s Status, success bool, t Thresholds) (Status, error) {

	newStatus := Status{
		State:     model.StateUp,
		Failures:  0,
		Successes: 0,
	}

	var err error

	if !success {
		if s.State == model.StateUnknown && s.Failures > t.Failure {
			newStatus = Status{}
			err = fmt.Errorf("Invalid State: %w", ErrInvalidUnknownStatus)
		}

		if s.Failures >= t.Failure {
			newStatus.State = model.StateDown
			newStatus.Failures = s.Failures + 1
			newStatus.Successes = 0
		} else if s.Failures < t.Failure {
			newStatus.State = s.State
			newStatus.Failures = s.Failures + 1
			newStatus.Successes = 0
		}
	}

	if success {
		if s.State == model.StateDown {
			err = fmt.Errorf("Invalid State: %w", ErrInvalidDownStatus)
		}

		if s.Successes >= t.Recovery {
			newStatus.State = model.StateUp
			newStatus.Failures = 0
			newStatus.Successes = s.Successes + 1
		} else if s.Successes < t.Recovery {
			newStatus.State = s.State
			newStatus.Failures = 0
			newStatus.Successes = s.Successes + 1
		}
	}

	return newStatus, err
}

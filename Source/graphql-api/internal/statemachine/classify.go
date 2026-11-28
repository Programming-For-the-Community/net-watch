package statemachine

import (
	"fmt"

	"github.com/Programming-For-the-Community/net-watch/graphql-api/internal/model"
)

// Classify reports the event that occurred between the previous status and the current status.
func Classify(prev, current Status) (Event, error) {

	var event = EventNone
	var err error

	if err = prev.Validate(); err != nil {
		return event, err
	}

	if err = current.Validate(); err != nil {
		return event, err
	}

	if prev.State != model.StateUnknown && current.State == model.StateUnknown {
		err = fmt.Errorf("invalid state transition: %w", ErrImpossibleTransition)
		return event, err
	}

	if prev.State != current.State && (current.State == model.StatePaused || prev.State == model.StatePaused) {
		err = fmt.Errorf("invalid state transition: %w", ErrImpossibleTransition)
		return event, err
	}

	if prev.State != model.StateDown && current.State == model.StateDown {
		event = EventWentDown
	}

	if prev.State == model.StateDown && current.State == model.StateUp {
		event = EventRecovered
	}

	if prev.State == model.StateUnknown && current.State == model.StateUp {
		event = EventCameUp
	}

	return event, err
}

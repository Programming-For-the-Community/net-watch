package statemachine

import "github.com/Programming-For-the-Community/net-watch/graphql-api/internal/model"

// Classify reports the event that occurred between the previous status and the current status.
func Classify(prev, current Status) (Event, error) {

	var event Event

	if prev.State != model.StateDown && current.State == model.StateDown {
		event = EventWentDown
	}

	if prev.State == model.StateDown && current.State == model.StateUp {
		event = EventRecovered
	}

	if prev.State == model.StateUnknown && current.State == model.StateUp {
		event = EventCameUp
	}

	return event, nil
}

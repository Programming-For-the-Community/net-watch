package statemachine

// Event is what happened during Next, classifies the Status transitions
type Event int

// Possible values of Event
const (
	EventNone      Event = 0
	EventWentDown  Event = 1
	EventCameUp    Event = 2
	EventRecovered Event = 3
)

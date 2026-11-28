package statemachine

import "fmt"

// Event is what happened during Next, classifies the Status transitions
type Event int

// Possible values of Event
const (
	EventNone      Event = 0
	EventWentDown  Event = 1
	EventCameUp    Event = 2
	EventRecovered Event = 3
)

func (e Event) String() string {
	switch e {
	case EventNone:
		return "EventNone"
	case EventWentDown:
		return "EventWentDown"
	case EventCameUp:
		return "EventCameUp"
	case EventRecovered:
		return "EventRecovered"
	default:
		return fmt.Sprintf("Event(%d)", int(e))
	}
}

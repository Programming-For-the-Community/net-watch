package statemachine

import (
	"testing"

	"github.com/Programming-For-the-Community/net-watch/graphql-api/internal/model"
)

// Classify reports what a single step of Next did, so callers know whether
// to open or resolve an incident and whether to send a notification.
//
// Suggested API (classify.go):
//
//	type Event int
//
//	const (
//		EventNone      Event = iota // no state change
//		EventWentDown               // tripped the failure threshold: open incident, notify
//		EventRecovered              // Down → Up: resolve incident, notify
//		EventCameUp                 // Unknown → Up: first decision, no incident
//	)
//
//	func Classify(prev, next Status) (Event, error)
//
// Classify only looks at State; counters never affect the result. Pausing
// and resuming are user actions handled by Pause/Resume, not by Classify.

// ----------- Transitions that matter -----------
// [ ] TestClassify_UpToDownReturnsWentDown
func TestClassify_UpToDownReturnsWentDown(t *testing.T) {
	prev := Status{State: model.StateUp}
	current := Status{State: model.StateDown}

	event, err := Classify(prev, current)

	if err != nil {
		t.Errorf("Classify(%v, %v) error = %v, want nil", prev, current, err)
	}

	if event != EventWentDown {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventWentDown)
	}
}

// [ ] TestClassify_UnknownToDownReturnsWentDown
func TestClassify_UnknownToDownReturnsWentDown(t *testing.T) {
	prev := Status{State: model.StateUnknown}
	current := Status{State: model.StateDown}

	event, err := Classify(prev, current)

	if err != nil {
		t.Errorf("Classify(%v, %v) error = %v, want nil", prev, current, err)
	}

	if event != EventWentDown {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventWentDown)
	}
}

// [ ] TestClassify_DownToUpReturnsRecovered
func TestClassify_DownToUpReturnsRecovered(t *testing.T) {
	prev := Status{State: model.StateDown}
	current := Status{State: model.StateUp}

	event, err := Classify(prev, current)

	if err != nil {
		t.Errorf("Classify(%v, %v) error = %v, want nil", prev, current, err)
	}

	if event != EventRecovered {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventRecovered)
	}
}

// [ ] TestClassify_UnknownToUpReturnsCameUp
func TestClassify_UnknownToUpReturnsCameUp(t *testing.T) {
	prev := Status{State: model.StateUnknown}
	current := Status{State: model.StateUp}

	event, err := Classify(prev, current)

	if err != nil {
		t.Errorf("Classify(%v, %v) error = %v, want nil", prev, current, err)
	}

	if event != EventCameUp {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventCameUp)
	}
}

// ----------- No transition -----------
// [ ] TestClassify_UpToUpReturnsNone
// [ ] TestClassify_DownToDownReturnsNone
// [ ] TestClassify_UnknownToUnknownReturnsNone
// [ ] TestClassify_PausedToPausedReturnsNone
// [ ] TestClassify_CounterOnlyChangeReturnsNone
//       e.g. Up{Failures: 1} → Up{Failures: 2}
// [ ] TestClassify_DownBuildingRecoverySuccessesReturnsNone
//       e.g. Down{Failures: 7} → Down{Successes: 1}

// ----------- Changes Next never makes -----------
// DECIDE: return an error (e.g. ErrImpossibleTransition) or EventNone.
// An error surfaces bugs in the caller; EventNone is more forgiving.
// [ ] TestClassify_UpToUnknownReturnsError
// [ ] TestClassify_DownToUnknownReturnsError
// [ ] TestClassify_UpToPausedReturnsError
// [ ] TestClassify_DownToPausedReturnsError
// [ ] TestClassify_UnknownToPausedReturnsError
// [ ] TestClassify_PausedToUpReturnsError
// [ ] TestClassify_PausedToDownReturnsError
// [ ] TestClassify_PausedToUnknownReturnsError

// ----------- Validation -----------
// [ ] TestClassify_EmptyPrevStateReturnsError
// [ ] TestClassify_EmptyNextStateReturnsError
// [ ] TestClassify_UnrecognisedPrevStateReturnsError
// [ ] TestClassify_UnrecognisedNextStateReturnsError
// [ ] TestClassify_ErrorReturnsEventNone
//       whatever the error, the Event returned alongside it is EventNone

// ----------- Full grid -----------
// [ ] TestClassify_AllStatePairs
//       table-driven over all 4×4 (prev, next) state pairs, so no
//       combination is missed; each row is either an Event or an error

// ----------- Working with Next -----------
// [ ] TestClassify_ConsecutiveFailuresFromUpReportWentDownExactlyOnce
//       run Next for t.Failure failures from Up, Classify each step:
//       EventNone until the last step, which is EventWentDown
// [ ] TestClassify_ConsecutiveSuccessesFromDownReportRecoveredExactlyOnce
// [ ] TestClassify_FirstDecisionFromUnknownReportsCameUp
//       with Recovery: 1, Unknown + success → EventCameUp
// [ ] TestClassify_FlappingReportsOnlyNone
//       alternating results with thresholds of 2 never produce an event
// [ ] TestClassify_PausedStepReportsNone
//       Next on a paused status returns it unchanged → EventNone

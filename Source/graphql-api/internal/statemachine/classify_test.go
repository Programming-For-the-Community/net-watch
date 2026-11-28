package statemachine

import (
	"errors"

	"github.com/Programming-For-the-Community/net-watch/graphql-api/internal/model"

	"testing"
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
func TestClassify_UpToUpReturnsNone(t *testing.T) {
	prev := Status{State: model.StateUp}
	current := Status{State: model.StateUp}

	event, err := Classify(prev, current)

	if err != nil {
		t.Errorf("Classify(%v, %v) error = %v, want nil", prev, current, err)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_DownToDownReturnsNone
func TestClassify_DownToDownReturnsNone(t *testing.T) {
	prev := Status{State: model.StateDown}
	current := Status{State: model.StateDown}

	event, err := Classify(prev, current)

	if err != nil {
		t.Errorf("Classify(%v, %v) error = %v, want nil", prev, current, err)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_UnknownToUnknownReturnsNone
func TestClassify_UnknownToUnknownReturnsNone(t *testing.T) {
	prev := Status{State: model.StateUnknown}
	current := Status{State: model.StateUnknown}

	event, err := Classify(prev, current)

	if err != nil {
		t.Errorf("Classify(%v, %v) error = %v, want nil", prev, current, err)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_PausedToPausedReturnsNone
func TestClassify_PausedToPausedReturnsNone(t *testing.T) {
	prev := Status{State: model.StatePaused}
	current := Status{State: model.StatePaused}

	event, err := Classify(prev, current)

	if err != nil {
		t.Errorf("Classify(%v, %v) error = %v, want nil", prev, current, err)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_CounterOnlyChangeReturnsNone
//
//	e.g. Up{Failures: 1} → Up{Failures: 2}
func TestClassify_CounterOnlyChangesReturnsNone(t *testing.T) {
	prev := Status{State: model.StateUp, Failures: 1, Successes: 0}
	current := Status{State: model.StateUp, Failures: 2, Successes: 0}

	event, err := Classify(prev, current)

	if err != nil {
		t.Errorf("Classify(%v, %v) error = %v, want nil", prev, current, err)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_DownBuildingRecoverySuccessesReturnsNone
//
//	e.g. Down{Failures: 7} → Down{Successes: 1}
func TestClassify_DownBuildingRecoverySuccessesReturnsNone(t *testing.T) {
	prev := Status{State: model.StateDown, Failures: 7, Successes: 0}
	current := Status{State: model.StateDown, Failures: 0, Successes: 1}

	event, err := Classify(prev, current)

	if err != nil {
		t.Errorf("Classify(%v, %v) error = %v, want nil", prev, current, err)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// ----------- Changes Next never makes -----------
// Return both Error AND EventNone for changes Next never makes since some sort
// of Event must be returned, nil is not an option
// [ ] TestClassify_UpToUnknownReturnsErrorAndEventNone
func TestClassify_UpToUnknownReturnsErrorAndEventNone(t *testing.T) {
	prev := Status{State: model.StateUp}
	current := Status{State: model.StateUnknown}

	event, err := Classify(prev, current)

	if !errors.Is(err, ErrImpossibleTransition) {
		t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, ErrImpossibleTransition)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_DownToUnknownReturnsErrorAndEventNone
func TestClassify_DownToUnknownReturnsErrorAndEventNone(t *testing.T) {
	prev := Status{State: model.StateDown}
	current := Status{State: model.StateUnknown}

	event, err := Classify(prev, current)

	if !errors.Is(err, ErrImpossibleTransition) {
		t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, ErrImpossibleTransition)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_UpToPausedReturnsErrorAndEventNone
func TestClassify_UpToPausedReturnsErrorAndEventNone(t *testing.T) {
	prev := Status{State: model.StateUp}
	current := Status{State: model.StatePaused}

	event, err := Classify(prev, current)

	if !errors.Is(err, ErrImpossibleTransition) {
		t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, ErrImpossibleTransition)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_DownToPausedReturnsErrorAndEventNone
func TestClassify_DownToPausedReturnsErrorAndEventNone(t *testing.T) {
	prev := Status{State: model.StateDown}
	current := Status{State: model.StatePaused}

	event, err := Classify(prev, current)

	if !errors.Is(err, ErrImpossibleTransition) {
		t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, ErrImpossibleTransition)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_UnknownToPausedReturnsErrorAndEventNone
func TestClassify_UnknownToPausedReturnsErrorAndEventNone(t *testing.T) {
	prev := Status{State: model.StateUnknown}
	current := Status{State: model.StatePaused}

	event, err := Classify(prev, current)

	if !errors.Is(err, ErrImpossibleTransition) {
		t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, ErrImpossibleTransition)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_PausedToUpReturnsErrorAndEventNone
func TestClassify_PausedToUpReturnsErrorAndEventNone(t *testing.T) {
	prev := Status{State: model.StatePaused}
	current := Status{State: model.StateUp}

	event, err := Classify(prev, current)

	if !errors.Is(err, ErrImpossibleTransition) {
		t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, ErrImpossibleTransition)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_PausedToDownReturnsErrorAndEventNone
func TestClassify_PausedToDownReturnsErrorAndEventNone(t *testing.T) {
	prev := Status{State: model.StatePaused}
	current := Status{State: model.StateDown}

	event, err := Classify(prev, current)

	if !errors.Is(err, ErrImpossibleTransition) {
		t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, ErrImpossibleTransition)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_PausedToUnknownReturnsErrorAndEventNone
func TestClassify_PausedToUnknownReturnsErrorAndEventNone(t *testing.T) {
	prev := Status{State: model.StatePaused}
	current := Status{State: model.StateUnknown}

	event, err := Classify(prev, current)

	if !errors.Is(err, ErrImpossibleTransition) {
		t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, ErrImpossibleTransition)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// ----------- Validation -----------
// [ ] TestClassify_EmptyPrevStateReturnsError
func TestClassify_EmptyPrevStateReturnsError(t *testing.T) {
	prev := Status{}
	current := Status{State: model.StateUnknown}

	event, err := Classify(prev, current)

	if !errors.Is(err, ErrStatusMissingStateField) {
		t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, ErrStatusMissingStateField)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_EmptyCurrentStateReturnsError
func TestClassify_EmptyCurrentStateReturnsError(t *testing.T) {
	prev := Status{State: model.StateUnknown}
	current := Status{}

	event, err := Classify(prev, current)

	if !errors.Is(err, ErrStatusMissingStateField) {
		t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, ErrStatusMissingStateField)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_UnrecognisedPrevStateReturnsError
func TestClassify_UnrecognisedPrevStateReturnsError(t *testing.T) {
	prev := Status{State: "some_state"}
	current := Status{State: model.StateUnknown}

	event, err := Classify(prev, current)

	if !errors.Is(err, ErrStatusUnknownState) {
		t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, ErrStatusUnknownState)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// [ ] TestClassify_UnrecognisedCurrentStateReturnsError
func TestClassify_UnrecognisedCurrentStateReturnsError(t *testing.T) {
	prev := Status{State: model.StateUnknown}
	current := Status{State: "some_state"}

	event, err := Classify(prev, current)

	if !errors.Is(err, ErrStatusUnknownState) {
		t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, ErrStatusUnknownState)
	}

	if event != EventNone {
		t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
	}
}

// ----------- Full grid -----------
// [ ] TestClassify_AllStatePairs
//
//	table-driven over all 4×4 (prev, next) state pairs, so no
//	combination is missed; each row is either an Event or an error
func TestClassify_AllStatePairs(t *testing.T) {
	tests := []struct {
		prev, current model.CheckState
		wantEvent     Event
		wantErr       error
	}{
		{model.StateUp, model.StateUp, EventNone, nil},
		{model.StateUp, model.StateDown, EventWentDown, nil},
		{model.StateUp, model.StateUnknown, EventNone, ErrImpossibleTransition},
		{model.StateUp, model.StatePaused, EventNone, ErrImpossibleTransition},

		{model.StateDown, model.StateUp, EventRecovered, nil},
		{model.StateDown, model.StateDown, EventNone, nil},
		{model.StateDown, model.StateUnknown, EventNone, ErrImpossibleTransition},
		{model.StateDown, model.StatePaused, EventNone, ErrImpossibleTransition},

		{model.StateUnknown, model.StateUp, EventCameUp, nil},
		{model.StateUnknown, model.StateDown, EventWentDown, nil},
		{model.StateUnknown, model.StateUnknown, EventNone, nil},
		{model.StateUnknown, model.StatePaused, EventNone, ErrImpossibleTransition},

		{model.StatePaused, model.StateUp, EventNone, ErrImpossibleTransition},
		{model.StatePaused, model.StateDown, EventNone, ErrImpossibleTransition},
		{model.StatePaused, model.StateUnknown, EventNone, ErrImpossibleTransition},
		{model.StatePaused, model.StatePaused, EventNone, nil},
	}

	for _, tc := range tests {
		t.Run(string(tc.prev)+"_to_"+string(tc.current), func(t *testing.T) {
			prev := Status{State: tc.prev}
			current := Status{State: tc.current}

			event, err := Classify(prev, current)

			if tc.wantErr == nil {
				if err != nil {
					t.Errorf("Classify(%v, %v) error = %v, want nil", prev, current, err)
				}
			} else if !errors.Is(err, tc.wantErr) {
				t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, tc.wantErr)
			}

			if event != tc.wantEvent {
				t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, tc.wantEvent)
			}
		})
	}
}

// ----------- Working with Next -----------

// stepAndClassify runs one Next step from s and classifies the transition,
// failing the test immediately if either call returns an error.
func stepAndClassify(t *testing.T, s Status, success bool, th Thresholds) (Status, Event) {
	t.Helper()
	next := mustNext(t, s, success, th)
	event, err := Classify(s, next)
	if err != nil {
		t.Fatalf("Classify(%+v, %+v) returned error: %v", s, next, err)
	}
	return next, event
}

// [ ] TestClassify_ConsecutiveFailuresFromUpReportWentDownExactlyOnce
//
//	run Next for t.Failure failures from Up, Classify each step:
//	EventNone until the last step, which is EventWentDown
func TestClassify_ConsecutiveFailuresFromUpReportWentDownExactlyOnce(t *testing.T) {
	th := Thresholds{Failure: 3, Recovery: 2}
	s := Status{State: model.StateUp}

	for i := 1; i <= th.Failure; i++ {
		var event Event
		s, event = stepAndClassify(t, s, false, th)

		want := EventNone
		if i == th.Failure {
			want = EventWentDown
		}
		if event != want {
			t.Errorf("failure %d: event = %v, want %v", i, event, want)
		}
	}

	assertStatus(t, s, Status{State: model.StateDown, Failures: th.Failure})
}

// [ ] TestClassify_ConsecutiveSuccessesFromDownReportRecoveredExactlyOnce
func TestClassify_ConsecutiveSuccessesFromDownReportRecoveredExactlyOnce(t *testing.T) {
	th := Thresholds{Failure: 2, Recovery: 3}
	s := Status{State: model.StateDown, Failures: th.Failure}

	for i := 1; i <= th.Recovery; i++ {
		var event Event
		s, event = stepAndClassify(t, s, true, th)

		want := EventNone
		if i == th.Recovery {
			want = EventRecovered
		}
		if event != want {
			t.Errorf("success %d: event = %v, want %v", i, event, want)
		}
	}

	assertStatus(t, s, Status{State: model.StateUp, Successes: th.Recovery})
}

// [ ] TestClassify_FirstDecisionFromUnknownReportsCameUp
//
//	with Recovery: 1, Unknown + success → EventCameUp
func TestClassify_FirstDecisionFromUnknownReportsCameUp(t *testing.T) {
	th := Thresholds{Failure: 3, Recovery: 1}

	s, event := stepAndClassify(t, Reset(), true, th)

	if event != EventCameUp {
		t.Errorf("event = %v, want %v", event, EventCameUp)
	}
	assertStatus(t, s, Status{State: model.StateUp, Successes: 1})
}

// [ ] TestClassify_FlappingReportsOnlyNone
//
//	alternating results with thresholds of 2 never produce an event
func TestClassify_FlappingReportsOnlyNone(t *testing.T) {
	th := Thresholds{Failure: 2, Recovery: 2}
	starts := []Status{
		{State: model.StateUnknown},
		{State: model.StateUp},
		{State: model.StateDown},
	}

	for _, start := range starts {
		t.Run(string(start.State), func(t *testing.T) {
			s := start
			for i := 0; i < 10; i++ {
				success := i%2 == 0
				var event Event
				s, event = stepAndClassify(t, s, success, th)
				if event != EventNone {
					t.Errorf("step %d (success=%v): event = %v, want %v", i, success, event, EventNone)
				}
			}
			if s.State != start.State {
				t.Errorf("state = %v, want %v", s.State, start.State)
			}
		})
	}
}

// [ ] TestClassify_PausedStepReportsNone
//
//	Next on a paused status returns it unchanged → EventNone
func TestClassify_PausedStepReportsNone(t *testing.T) {
	th := Thresholds{Failure: 1, Recovery: 1}
	paused := Status{State: model.StatePaused}

	for _, success := range []bool{true, false} {
		s, event := stepAndClassify(t, paused, success, th)
		if event != EventNone {
			t.Errorf("success=%v: event = %v, want %v", success, event, EventNone)
		}
		assertStatus(t, s, paused)
	}
}

// [ ] TestClassify_FailuresWhileDownReportOnlyNone
//
//	once Down, further failures must not re-report EventWentDown
func TestClassify_FailuresWhileDownReportOnlyNone(t *testing.T) {
	th := Thresholds{Failure: 2, Recovery: 2}
	s := Status{State: model.StateDown, Failures: th.Failure}

	for i := 1; i <= 5; i++ {
		var event Event
		s, event = stepAndClassify(t, s, false, th)
		if event != EventNone {
			t.Errorf("extra failure %d: event = %v, want %v", i, event, EventNone)
		}
	}

	assertStatus(t, s, Status{State: model.StateDown, Failures: th.Failure + 5})
}

// [ ] TestClassify_SuccessesWhileUpReportOnlyNone
//
//	once Up, further successes must not re-report EventRecovered or EventCameUp
func TestClassify_SuccessesWhileUpReportOnlyNone(t *testing.T) {
	th := Thresholds{Failure: 2, Recovery: 2}
	s := Status{State: model.StateUp, Successes: th.Recovery}

	for i := 1; i <= 5; i++ {
		var event Event
		s, event = stepAndClassify(t, s, true, th)
		if event != EventNone {
			t.Errorf("extra success %d: event = %v, want %v", i, event, EventNone)
		}
	}

	assertStatus(t, s, Status{State: model.StateUp, Successes: th.Recovery + 5})
}

// [ ] TestClassify_FullOutageCycleReportsEachEventOnce
//
//	Up → Down → Up via Next: exactly one EventWentDown then one EventRecovered
func TestClassify_FullOutageCycleReportsEachEventOnce(t *testing.T) {
	th := Thresholds{Failure: 3, Recovery: 2}
	s := Status{State: model.StateUp}

	// 5 failures (2 past the threshold), then 4 successes (2 past the threshold).
	results := []bool{false, false, false, false, false, true, true, true, true}
	var got []Event
	for _, success := range results {
		var event Event
		s, event = stepAndClassify(t, s, success, th)
		if event != EventNone {
			got = append(got, event)
		}
	}

	want := []Event{EventWentDown, EventRecovered}
	if len(got) != len(want) {
		t.Fatalf("events = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("events = %v, want %v", got, want)
			break
		}
	}
	assertStatus(t, s, Status{State: model.StateUp, Successes: 4})
}

// [ ] TestClassify_FirstDecisionFromUnknownReportsWentDown
//
//	reaching the failure threshold from Unknown → EventWentDown on that step only
func TestClassify_FirstDecisionFromUnknownReportsWentDown(t *testing.T) {
	th := Thresholds{Failure: 3, Recovery: 2}
	s := Reset()

	for i := 1; i <= th.Failure; i++ {
		var event Event
		s, event = stepAndClassify(t, s, false, th)

		want := EventNone
		if i == th.Failure {
			want = EventWentDown
		}
		if event != want {
			t.Errorf("failure %d: event = %v, want %v", i, event, want)
		}
	}

	assertStatus(t, s, Status{State: model.StateDown, Failures: th.Failure})
}

// [ ] TestClassify_InvalidCountersReturnError
//
//	Classify must surface Validate's counter errors for either input
func TestClassify_InvalidCountersReturnError(t *testing.T) {
	valid := Status{State: model.StateUp}
	bad := []struct {
		name    string
		status  Status
		wantErr error
	}{
		{"negative_failures", Status{State: model.StateUp, Failures: -1}, ErrStatusFailuresFieldNegative},
		{"negative_successes", Status{State: model.StateUp, Successes: -1}, ErrStatusSuccessesFieldNegative},
		{"both_counters_set", Status{State: model.StateUp, Failures: 1, Successes: 1}, ErrInvalidFailuresAndSuccessesCounts},
	}

	for _, tc := range bad {
		for _, side := range []string{"prev", "current"} {
			t.Run(side+"_"+tc.name, func(t *testing.T) {
				prev, current := valid, tc.status
				if side == "prev" {
					prev, current = tc.status, valid
				}

				event, err := Classify(prev, current)

				if !errors.Is(err, tc.wantErr) {
					t.Errorf("Classify(%v, %v) error = %v, want %v", prev, current, err, tc.wantErr)
				}
				if event != EventNone {
					t.Errorf("Classify(%v, %v) event = %v, want %v", prev, current, event, EventNone)
				}
			})
		}
	}
}

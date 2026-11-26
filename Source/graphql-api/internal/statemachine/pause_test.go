package statemachine

import (
	"testing"

	"github.com/Programming-For-the-Community/net-watch/graphql-api/internal/model"
)

// Pause and Resume change a status in response to a user action. Paused is a
// stored state (check_status.state = 'paused'); resuming starts the check
// over from Unknown so stale streaks don't carry across the pause.
//
// Suggested API (pause.go):
//
//	func Pause(s Status) (Status, error)
//	func Resume(s Status) (Status, error)
//
// Resume returns the same value as Reset (already in pause.go; its test is
// at the bottom of this file). The difference is that Resume checks it's
// being applied to a paused status.

// ----------- Pause -----------
// [ ] TestPause_UpStatusBecomesPaused
// [ ] TestPause_DownStatusBecomesPaused
// [ ] TestPause_UnknownStatusBecomesPaused
// [ ] TestPause_ClearsCounters
//       e.g. Up{Failures: 2} → Paused{0, 0}. Counters mean nothing while
//       paused and are reset on resume anyway.
// DECIDE: pausing an already-paused status is a no-op or an error.
// A no-op makes "pause" safe to repeat (e.g. from the AI assistant).
// [ ] TestPause_AlreadyPausedStatusIsUnchanged
//       (or TestPause_AlreadyPausedStatusReturnsError)
// [ ] TestPause_EmptyStateReturnsError
// [ ] TestPause_UnrecognisedStateReturnsError
// [ ] TestPause_NegativeCountersReturnError
//       Pause has no thresholds, so it can only check fields on their own

// ----------- Resume -----------
// [ ] TestResume_PausedStatusBecomesUnknown
// [ ] TestResume_ResetsCountersToZero
//       e.g. Paused{Failures: 3} → Unknown{0, 0}
// [ ] TestResume_MatchesReset
//       Resume(paused) == Reset()
// DECIDE: resuming a status that isn't paused is a no-op or an error.
// If it's a no-op, it must return s unchanged, NOT Reset(); otherwise
// "resume" on a Down check would silently wipe its state.
// [ ] TestResume_UpStatusIsUnchanged         (or ...ReturnsError)
// [ ] TestResume_DownStatusIsUnchanged       (or ...ReturnsError)
// [ ] TestResume_UnknownStatusIsUnchanged    (or ...ReturnsError)
// [ ] TestResume_EmptyStateReturnsError
// [ ] TestResume_UnrecognisedStateReturnsError

// ----------- Round trips -----------
// [ ] TestPauseResume_RoundTripFromAnyStateEndsUnknown
//       Pause then Resume from Up, Down and Unknown all end at Reset()
// [ ] TestPauseResume_NextAfterResumeStartsFromScratch
//       Down{Failures: 9} → Pause → Resume → Next(success) with
//       Recovery: 2 gives Unknown{Successes: 1}, not Up
// [ ] TestPause_NextWhilePausedIgnoresResults
//       Pause → several Next calls (both results) → still Paused{0, 0}

// ----------- Open questions for later pieces (not tests here) -----------
// - Pausing while Down with an open incident: resolve it, leave it open,
//   or mark it some other way? This belongs to the incident logic, but
//   decide it before writing those tests.
// - Should pausing/resuming produce events for notifications or the UI?
//   If so, add EventPaused/EventResumed and return them from Pause/Resume
//   rather than from Classify.

// ----------- Reset -----------

func TestReset_ReturnsUnknownWithZeroCounters(t *testing.T) {
	assertStatus(t, Reset(), Status{State: model.StateUnknown})
}

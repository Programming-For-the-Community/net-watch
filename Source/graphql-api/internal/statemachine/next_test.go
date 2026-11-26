package statemachine

import (
	"testing"

	"github.com/Programming-For-the-Community/net-watch/graphql-api/internal/model"
)

var defaultThresholds = Thresholds{Failure: 5, Recovery: 2}

// ----------- Input validation -----------

func TestNext_RequiresStateFieldInStatusStruct(t *testing.T) {
	assertNextError(t, Status{}, true, defaultThresholds, ErrStatusMissingStateField)
}

func TestNext_RejectsUnrecognisedState(t *testing.T) {
	assertNextError(t, Status{State: "banana"}, true, defaultThresholds, ErrStatusUnknownState)
}

func TestNext_RequiresNonNegativeFailuresInStatusStruct(t *testing.T) {
	assertNextError(t, Status{State: model.StateUp, Failures: -1}, true, defaultThresholds, ErrStatusFailuresFieldNegative)
}

func TestNext_RequiresNonNegativeSuccessesInStatusStruct(t *testing.T) {
	assertNextError(t, Status{State: model.StateUp, Successes: -1}, true, defaultThresholds, ErrStatusSuccessesFieldNegative)
}

func TestNext_RequiresEitherFailuresOrSuccessesBeZero(t *testing.T) {
	s := Status{State: model.StateUp, Failures: 2, Successes: 4}
	assertNextError(t, s, true, defaultThresholds, ErrInvalidFailuresAndSuccessesCounts)
}

func TestNext_RequiresNonNegativeFailureInThresholdsStruct(t *testing.T) {
	assertNextError(t, Status{State: model.StateUp}, true, Thresholds{Failure: -1, Recovery: 2}, ErrThresholdsFailureTooLow)
}

func TestNext_RequiresNonZeroFailureInThresholdsStruct(t *testing.T) {
	assertNextError(t, Status{State: model.StateUp}, true, Thresholds{Failure: 0, Recovery: 2}, ErrThresholdsFailureTooLow)
}

func TestNext_RequiresNonNegativeRecoveryInThresholdsStruct(t *testing.T) {
	assertNextError(t, Status{State: model.StateUp}, true, Thresholds{Failure: 5, Recovery: -1}, ErrThresholdsRecoveryTooLow)
}

func TestNext_RequiresNonZeroRecoveryInThresholdsStruct(t *testing.T) {
	assertNextError(t, Status{State: model.StateUp}, true, Thresholds{Failure: 5, Recovery: 0}, ErrThresholdsRecoveryTooLow)
}

// A streak that has reached its threshold would already have changed state,
// so these statuses can never come out of Next.

func TestNext_UnknownStatusWithFailuresAtFailureThresholdReturnsError(t *testing.T) {
	s := Status{State: model.StateUnknown, Failures: defaultThresholds.Failure}
	assertNextError(t, s, false, defaultThresholds, ErrInvalidUnknownFailuresStatus)
}

func TestNext_UnknownStatusWithSuccessesAtRecoveryThresholdReturnsError(t *testing.T) {
	s := Status{State: model.StateUnknown, Successes: defaultThresholds.Recovery}
	assertNextError(t, s, true, defaultThresholds, ErrInvalidUnknownSuccessesStatus)
}

func TestNext_DownStatusWithSuccessesAtRecoveryThresholdReturnsError(t *testing.T) {
	s := Status{State: model.StateDown, Successes: defaultThresholds.Recovery}
	assertNextError(t, s, true, defaultThresholds, ErrInvalidDownStatus)
}

func TestNext_UpStatusWithFailuresAtFailureThresholdReturnsError(t *testing.T) {
	s := Status{State: model.StateUp, Failures: defaultThresholds.Failure}
	assertNextError(t, s, false, defaultThresholds, ErrInvalidUpStatus)
}

// ----------- Unknown -----------

func TestNext_SuccessfulUnknownStatusBelowRecoveryThresholdStaysUnknown(t *testing.T) {
	got := mustNext(t, Status{State: model.StateUnknown}, true, defaultThresholds)
	assertStatus(t, got, Status{State: model.StateUnknown, Successes: 1})
}

func TestNext_SuccessfulUnknownStatusReachingRecoveryThresholdBecomesUp(t *testing.T) {
	s := Status{State: model.StateUnknown, Successes: defaultThresholds.Recovery - 1}
	got := mustNext(t, s, true, defaultThresholds)
	assertStatus(t, got, Status{State: model.StateUp, Successes: defaultThresholds.Recovery})
}

func TestNext_FailedUnknownStatusBelowFailureThresholdStaysUnknown(t *testing.T) {
	s := Status{State: model.StateUnknown, Failures: defaultThresholds.Failure - 2}
	got := mustNext(t, s, false, defaultThresholds)
	assertStatus(t, got, Status{State: model.StateUnknown, Failures: defaultThresholds.Failure - 1})
}

func TestNext_FailedUnknownStatusReachingFailureThresholdBecomesDown(t *testing.T) {
	s := Status{State: model.StateUnknown, Failures: defaultThresholds.Failure - 1}
	got := mustNext(t, s, false, defaultThresholds)
	assertStatus(t, got, Status{State: model.StateDown, Failures: defaultThresholds.Failure})
}

func TestNext_SuccessfulUnknownStatusResetsFailureCount(t *testing.T) {
	got := mustNext(t, Status{State: model.StateUnknown, Failures: 2}, true, defaultThresholds)
	assertStatus(t, got, Status{State: model.StateUnknown, Successes: 1})
}

func TestNext_FailedUnknownStatusResetsSuccessCount(t *testing.T) {
	got := mustNext(t, Status{State: model.StateUnknown, Successes: 1}, false, defaultThresholds)
	assertStatus(t, got, Status{State: model.StateUnknown, Failures: 1})
}

// ----------- Up -----------

func TestNext_SuccessfulUpStatusStaysUp(t *testing.T) {
	got := mustNext(t, Status{State: model.StateUp, Successes: 7}, true, defaultThresholds)
	assertStatus(t, got, Status{State: model.StateUp, Successes: 8})
}

func TestNext_FailedUpStatusBelowFailureThresholdStaysUp(t *testing.T) {
	s := Status{State: model.StateUp, Failures: defaultThresholds.Failure - 2}
	got := mustNext(t, s, false, defaultThresholds)
	assertStatus(t, got, Status{State: model.StateUp, Failures: defaultThresholds.Failure - 1})
}

func TestNext_FailedUpStatusReachingFailureThresholdBecomesDown(t *testing.T) {
	s := Status{State: model.StateUp, Failures: defaultThresholds.Failure - 1}
	got := mustNext(t, s, false, defaultThresholds)
	assertStatus(t, got, Status{State: model.StateDown, Failures: defaultThresholds.Failure})
}

func TestNext_SuccessResetsFailureCount(t *testing.T) {
	s := Status{State: model.StateUp, Failures: defaultThresholds.Failure - 1}
	got := mustNext(t, s, true, defaultThresholds)
	assertStatus(t, got, Status{State: model.StateUp, Successes: 1})
}

// ----------- Down -----------

func TestNext_FailedDownStatusStaysDown(t *testing.T) {
	got := mustNext(t, Status{State: model.StateDown, Failures: 9}, false, defaultThresholds)
	assertStatus(t, got, Status{State: model.StateDown, Failures: 10})
}

func TestNext_SuccessfulDownStatusBelowRecoveryThresholdStaysDown(t *testing.T) {
	th := Thresholds{Failure: 5, Recovery: 3}
	s := Status{State: model.StateDown, Successes: th.Recovery - 2}
	got := mustNext(t, s, true, th)
	assertStatus(t, got, Status{State: model.StateDown, Successes: th.Recovery - 1})
}

func TestNext_SuccessfulDownStatusReachingRecoveryThresholdBecomesUp(t *testing.T) {
	s := Status{State: model.StateDown, Successes: defaultThresholds.Recovery - 1}
	got := mustNext(t, s, true, defaultThresholds)
	assertStatus(t, got, Status{State: model.StateUp, Successes: defaultThresholds.Recovery})
}

func TestNext_FailureResetsSuccessCount(t *testing.T) {
	s := Status{State: model.StateDown, Successes: defaultThresholds.Recovery - 1}
	got := mustNext(t, s, false, defaultThresholds)
	assertStatus(t, got, Status{State: model.StateDown, Failures: 1})
}

// ----------- Sequences -----------

func TestNext_ConsecutiveFailuresGoDownExactlyOnFailureThreshold(t *testing.T) {
	s := Status{State: model.StateUp}
	for i := 1; i <= defaultThresholds.Failure; i++ {
		s = mustNext(t, s, false, defaultThresholds)
		want := model.StateUp
		if i == defaultThresholds.Failure {
			want = model.StateDown
		}
		if s.State != want {
			t.Fatalf("after %d failures: state = %q, want %q", i, s.State, want)
		}
	}
}

func TestNext_ConsecutiveSuccessesGoUpExactlyOnRecoveryThreshold(t *testing.T) {
	th := Thresholds{Failure: 5, Recovery: 4}
	s := Status{State: model.StateDown}
	for i := 1; i <= th.Recovery; i++ {
		s = mustNext(t, s, true, th)
		want := model.StateDown
		if i == th.Recovery {
			want = model.StateUp
		}
		if s.State != want {
			t.Fatalf("after %d successes: state = %q, want %q", i, s.State, want)
		}
	}
}

// Thresholds of 2 mean that a single counter that failed to reset would
// trip the threshold on the second iteration.

func TestNext_FlappingWhileUpNeverReachesFailureThreshold(t *testing.T) {
	th := Thresholds{Failure: 2, Recovery: 2}
	s := Status{State: model.StateUp}
	for range 10 {
		s = mustNext(t, s, false, th)
		assertStatus(t, s, Status{State: model.StateUp, Failures: 1})
		s = mustNext(t, s, true, th)
		assertStatus(t, s, Status{State: model.StateUp, Successes: 1})
	}
}

func TestNext_FlappingWhileDownNeverReachesRecoveryThreshold(t *testing.T) {
	th := Thresholds{Failure: 2, Recovery: 2}
	s := Status{State: model.StateDown}
	for range 10 {
		s = mustNext(t, s, true, th)
		assertStatus(t, s, Status{State: model.StateDown, Successes: 1})
		s = mustNext(t, s, false, th)
		assertStatus(t, s, Status{State: model.StateDown, Failures: 1})
	}
}

// ----------- Paused -----------

func TestNext_SuccessfulPausedStatusStaysPausedWithCountersUnchanged(t *testing.T) {
	s := Status{State: model.StatePaused, Failures: 3}
	got := mustNext(t, s, true, defaultThresholds)
	assertStatus(t, got, s)
}

func TestNext_FailedPausedStatusStaysPausedWithCountersUnchanged(t *testing.T) {
	s := Status{State: model.StatePaused, Failures: 3}
	got := mustNext(t, s, false, defaultThresholds)
	assertStatus(t, got, s)
}

// Thresholds can change while a check is paused, so a paused streak above
// the current threshold is still valid.
func TestNext_PausedStatusIsNotComparedAgainstThresholds(t *testing.T) {
	s := Status{State: model.StatePaused, Failures: defaultThresholds.Failure + 10}
	got := mustNext(t, s, false, defaultThresholds)
	assertStatus(t, got, s)
}

func TestNext_PausedStatusWithBothCountersSetReturnsError(t *testing.T) {
	s := Status{State: model.StatePaused, Failures: 2, Successes: 3}
	assertNextError(t, s, true, defaultThresholds, ErrInvalidFailuresAndSuccessesCounts)
}

// ----------- Threshold edge cases -----------

func TestNext_FailureThresholdOfOneGoesDownOnFirstFailure(t *testing.T) {
	th := Thresholds{Failure: 1, Recovery: 2}
	got := mustNext(t, Status{State: model.StateUp, Successes: 1}, false, th)
	assertStatus(t, got, Status{State: model.StateDown, Failures: 1})
}

func TestNext_FailureThresholdOfOneGoesDownOnFirstFailureFromUnknown(t *testing.T) {
	th := Thresholds{Failure: 1, Recovery: 2}
	got := mustNext(t, Status{State: model.StateUnknown}, false, th)
	assertStatus(t, got, Status{State: model.StateDown, Failures: 1})
}

func TestNext_RecoveryThresholdOfOneGoesUpOnFirstSuccess(t *testing.T) {
	th := Thresholds{Failure: 5, Recovery: 1}
	got := mustNext(t, Status{State: model.StateDown, Failures: 3}, true, th)
	assertStatus(t, got, Status{State: model.StateUp, Successes: 1})
}

func TestNext_RecoveryThresholdOfOneGoesUpOnFirstSuccessFromUnknown(t *testing.T) {
	th := Thresholds{Failure: 5, Recovery: 1}
	got := mustNext(t, Status{State: model.StateUnknown}, true, th)
	assertStatus(t, got, Status{State: model.StateUp, Successes: 1})
}

// ----------- Errors -----------

func TestNext_ValidTransitionsReturnNilError(t *testing.T) {
	inputs := []Status{
		{State: model.StateUnknown},
		{State: model.StateUp, Failures: defaultThresholds.Failure - 1},
		{State: model.StateDown, Successes: defaultThresholds.Recovery - 1},
		{State: model.StatePaused},
	}
	for _, s := range inputs {
		for _, success := range []bool{true, false} {
			if _, err := Next(s, success, defaultThresholds); err != nil {
				t.Errorf("Next(%+v, %v) error = %v, want nil", s, success, err)
			}
		}
	}
}

package statemachine

import (
	"errors"
	"github.com/Programming-For-the-Community/net-watch/graphql-api/internal/model"
	"testing"
)

// Test list:
// [ ] unknown + success → up
func TestNext_UnknownWithSuccessBecomesUp(t *testing.T) {
	status, _ := Next(Status{State: model.StateUnknown}, true, Thresholds{Failure: 3, Recovery: 2})

	if status.State != model.StateUp {
		t.Errorf("state = %q, want %q", status.State, model.StateUp)
	}
}

// [ ] unknown + failure reaching threshold → down, transition reported
func TestNext_FailedUnknownStatusWithFailuresEqualToThresholdBecomesDown(t *testing.T) {
	testThreshold := Thresholds{Failure: 5, Recovery: 2}

	status, _ := Next(Status{State: model.StateUnknown, Failures: testThreshold.Failure}, false, testThreshold)

	if status.State != model.StateDown {
		t.Errorf("state = %q, want %q", status.State, model.StateDown)
	}

	if status.Failures != testThreshold.Failure+1 {
		t.Errorf("failures = %d, want %d", status.Failures, testThreshold.Failure+1)
	}
}

// [ ] unknown + failure below threshold → stays unknown, failures counted
func TestNext_FailedUnknownStatusWithLessFailuresThanThresholdStaysUnknown(t *testing.T) {
	testThreshold := Thresholds{Failure: 5, Recovery: 2}
	status, _ := Next(Status{State: model.StateUnknown, Failures: testThreshold.Failure - 1}, false, testThreshold)

	if status.State != model.StateUnknown {
		t.Errorf("state = %q, want %q", status.State, model.StateUnknown)
	}

	if status.Failures != testThreshold.Failure {
		t.Errorf("failures = %d, want %d", status.Failures, testThreshold.Failure)
	}
}

// [ ] unknown + failure above threshold → error
func TestNext_FailedUnknownStatusWithMoreFailuresThanThresholdFailuresThrowsError(t *testing.T) {
	testThreshold := Thresholds{Failure: 5, Recovery: 2}
	_, err := Next(Status{State: model.StateUnknown, Failures: testThreshold.Failure + 1}, false, testThreshold)

	if !errors.Is(err, ErrInvalidUnknownStatus) {
		t.Errorf("error = %v, want %v", err, ErrInvalidUnknownStatus)
	}
}

// [ ] up + failure below failure threshold → stays up, failures counted
func TestNext_FailedUpStatusWithLessFailuresThanThresholdStaysUp(t *testing.T) {
	testThreshold := Thresholds{Failure: 5, Recovery: 2}
	status, _ := Next(Status{State: model.StateUp, Failures: testThreshold.Failure - 1}, false, testThreshold)

	if status.State != model.StateUp {
		t.Errorf("state = %q, want %q", status.State, model.StateUp)
	}

	if status.Failures != testThreshold.Failure {
		t.Errorf("failures = %d, want %d", status.Failures, testThreshold.Failure)
	}
}

// [ ] up + failure reaching the threshold → down, transition reported
func TestNext_FailedUpStatusWithFailuresEqualToThresholdBecomesDown(t *testing.T) {
	testThreshold := Thresholds{Failure: 5, Recovery: 2}
	status, _ := Next(Status{State: model.StateUp, Failures: testThreshold.Failure}, false, testThreshold)

	if status.State != model.StateDown {
		t.Errorf("state = %q, want %q", status.State, model.StateDown)
	}

	if status.Failures != testThreshold.Failure+1 {
		t.Errorf("failures = %d, want %d", status.Failures, testThreshold.Failure+1)
	}
}

// [ ] down + success below recovery threshold → stays down
func TestNext_SuccessfulDownStatusWithLessSuccessesThanRecoveryThresholdStaysDown(t *testing.T) {
	testThreshold := Thresholds{Failure: 5, Recovery: 2}
	status, _ := Next(Status{State: model.StateDown, Successes: testThreshold.Recovery - 1}, true, testThreshold)

	if status.State != model.StateDown {
		t.Errorf("state = %q, want %q", status.State, model.StateDown)
	}

	if status.Successes != testThreshold.Recovery {
		t.Errorf("successes = %d, want %d", status.Successes, testThreshold.Recovery)
	}
}

// [ ] down + success reaching recovery threshold → up, transition reported
func TestNext_SuccessfulDownStatusWithSuccessesEqualToRecoveryThresholdBecomesUp(t *testing.T) {
	testThreshold := Thresholds{Failure: 5, Recovery: 2}
	status, _ := Next(Status{State: model.StateDown, Successes: testThreshold.Recovery}, true, testThreshold)

	if status.State != model.StateUp {
		t.Errorf("state = %q, want %q", status.State, model.StateUp)
	}

	if status.Successes != testThreshold.Recovery+1 {
		t.Errorf("successes = %d, want %d", status.Successes, testThreshold.Recovery+1)
	}
}

// [ ] donw + success above recovery threshold → error
func TestNext_SuccessfulDownStatusWithMoreSuccessesThanRecoveryThresholdThrowsError(t *testing.T) {
	testThreshold := Thresholds{Failure: 5, Recovery: 2}
	_, err := Next(Status{State: model.StateDown, Successes: testThreshold.Recovery + 1}, true, testThreshold)

	if !errors.Is(err, ErrInvalidDownStatus) {
		t.Errorf("error = %v, want %v", err, ErrInvalidDownStatus)
	}
}

// [ ] a success resets the failure count (and vice versa)
func TestNext_SuccessResetsFailureCount(t *testing.T) {
	testThreshold := Thresholds{Failure: 5, Recovery: 2}
	status, _ := Next(Status{State: model.StateUp, Failures: testThreshold.Failure - 1}, true, testThreshold)

	if status.Failures != 0 {
		t.Errorf("failures = %d, want %d", status.Failures, 0)
	}
}

func TestNext_FailureResetsSuccessCount(t *testing.T) {
	testThreshold := Thresholds{Failure: 5, Recovery: 2}
	status, _ := Next(Status{State: model.StateDown, Successes: testThreshold.Recovery - 1}, false, testThreshold)

	if status.Successes != 0 {
		t.Errorf("successes = %d, want %d", status.Successes, 0)
	}
}

// [ ] failure, success, failure (flapping) never trips the threshold
// [ ] paused → ignores results

// --- Steady states (no transition) ---
// [ ] TestNext_SuccessfulUpStatusStaysUp
// [ ] TestNext_FailedDownStatusStaysDown
// [ ] TestNext_FlappingWhileDownNeverReachesRecoveryThreshold

// --- Paused ---
// [ ] TestNext_SuccessfulPausedStatusStaysPausedWithCountersUnchanged
// [ ] TestNext_FailedPausedStatusStaysPausedWithCountersUnchanged

// --- Threshold edge cases ---
// [ ] TestNext_FailureThresholdOfOneGoesDownOnFirstFailure
// [ ] TestNext_RecoveryThresholdOfOneGoesUpOnFirstSuccess
// [ ] TestNext_FailedUpStatusWithMoreFailuresThanThresholdThrowsError

// --- Errors ---
// [ ] TestNext_ValidTransitionsReturnNilError
// [ ] TestNext_SuccessfulDownStatusWithLessSuccessesThanRecoveryThresholdReturnsNoError

// --- Input validation ---
// [ ] TestNext_EmptyStateTreatedAsUnknown
// [ ] TestNext_InvalidStateReturnsError
// [ ] TestNext_NegativeFailuresReturnsError
// [ ] TestNext_NegativeSuccessesReturnsError

// --- Thresholds.Validate() ---
// [ ] TestThresholdsValidate_FailureBelowOneReturnsError
// [ ] TestThresholdsValidate_RecoveryBelowOneReturnsError
// [ ] TestThresholdsValidate_ValidThresholdsReturnNil

// --- Later: Reset (used when the network changes) ---
// [ ] TestReset_ReturnsUnknownWithZeroCounters

// Decisions to make before writing the tests above:
//  - Does a threshold of N go down on the Nth failure or the (N+1)th? (current tests expect N+1)
//  - On error, return the zero Status or the input Status unchanged?

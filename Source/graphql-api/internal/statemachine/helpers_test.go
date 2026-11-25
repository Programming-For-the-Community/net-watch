package statemachine

import (
	"errors"
	"testing"
)

// mustNext calls Next and fails the test immediately if it returns an error.
func mustNext(t *testing.T, s Status, success bool, th Thresholds) Status {
	t.Helper()
	got, err := Next(s, success, th)
	if err != nil {
		t.Fatalf("Next(%+v, %v, %+v) returned error: %v", s, success, th, err)
	}
	return got
}

// assertNextError checks that Next rejects the input with wantErr and
// returns an empty Status.
func assertNextError(t *testing.T, s Status, success bool, th Thresholds, wantErr error) {
	t.Helper()
	got, err := Next(s, success, th)
	if !errors.Is(err, wantErr) {
		t.Errorf("Next(%+v, %v, %+v) error = %v, want %v", s, success, th, err, wantErr)
	}
	if got != (Status{}) {
		t.Errorf("Next(%+v, %v, %+v) status = %+v, want empty Status", s, success, th, got)
	}
}

func assertStatus(t *testing.T, got, want Status) {
	t.Helper()
	if got != want {
		t.Errorf("status = %+v, want %+v", got, want)
	}
}

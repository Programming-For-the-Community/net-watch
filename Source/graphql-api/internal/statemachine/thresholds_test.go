package statemachine

import (
	"errors"
	"testing"
)

func TestThresholdsValidate_FailureBelowOneReturnsError(t *testing.T) {
	err := Thresholds{Failure: 0, Recovery: 2}.Validate()
	if !errors.Is(err, ErrThresholdsFailureTooLow) {
		t.Errorf("error = %v, want %v", err, ErrThresholdsFailureTooLow)
	}
}

func TestThresholdsValidate_RecoveryBelowOneReturnsError(t *testing.T) {
	err := Thresholds{Failure: 5, Recovery: 0}.Validate()
	if !errors.Is(err, ErrThresholdsRecoveryTooLow) {
		t.Errorf("error = %v, want %v", err, ErrThresholdsRecoveryTooLow)
	}
}

func TestThresholdsValidate_ValidThresholdsReturnNil(t *testing.T) {
	if err := defaultThresholds.Validate(); err != nil {
		t.Errorf("error = %v, want nil", err)
	}
}

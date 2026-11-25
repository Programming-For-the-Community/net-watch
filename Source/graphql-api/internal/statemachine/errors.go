package statemachine

import "errors"

// Status field errors.
var (
	ErrStatusMissingStateField           = errors.New("status must have a state")
	ErrStatusUnknownState                = errors.New("status has an unrecognised state")
	ErrStatusFailuresFieldNegative       = errors.New("status failures must not be negative")
	ErrStatusSuccessesFieldNegative      = errors.New("status successes must not be negative")
	ErrInvalidFailuresAndSuccessesCounts = errors.New("status cannot have both failures and successes; one must be 0")
)

// Thresholds field errors.
var (
	ErrThresholdsFailureTooLow  = errors.New("failure threshold must be at least 1")
	ErrThresholdsRecoveryTooLow = errors.New("recovery threshold must be at least 1")
)

// Status/threshold combinations that Next can never produce.
var (
	ErrInvalidUnknownFailuresStatus  = errors.New("unknown status cannot have failures at or above the failure threshold")
	ErrInvalidUnknownSuccessesStatus = errors.New("unknown status cannot have successes at or above the recovery threshold")
	ErrInvalidDownStatus             = errors.New("down status cannot have successes at or above the recovery threshold")
	ErrInvalidUpStatus               = errors.New("up status cannot have failures at or above the failure threshold")
)

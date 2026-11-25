package statemachine

import "errors"

var ErrInvalidUnknownStatus = errors.New("Status cannot be unknown with failures exceeding threshold failures")

var ErrInvalidDownStatus = errors.New("Status cannot be down with successes exceeding threshold recoveries")

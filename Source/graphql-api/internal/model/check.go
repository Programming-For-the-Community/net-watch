// Package model holds the domain types shared across the graphql-api.
package model

// CheckState is the health state of a check.
type CheckState string

// Possible values of CheckState.
const (
	StateUnknown CheckState = "unknown" // not enough results yet to decide
	StateUp      CheckState = "up"
	StateDown    CheckState = "down"
	StatePaused  CheckState = "paused" // results are ignored until resumed
)

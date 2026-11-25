package model

type CheckState string

const (
	StateUnknown CheckState = "unknown"
	StateUp      CheckState = "up"
	StateDown    CheckState = "down"
	StatePaused  CheckState = "paused"
)

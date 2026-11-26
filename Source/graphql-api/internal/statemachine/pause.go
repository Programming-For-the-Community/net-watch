package statemachine

import (
	"github.com/Programming-For-the-Community/net-watch/graphql-api/internal/model"
)

// Reset returns a status with no history: unknown, with both counters at 0.
// Use it when a check is created, resumed after a pause, or reconfigured.
func Reset() Status {
	return Status{State: model.StateUnknown}
}

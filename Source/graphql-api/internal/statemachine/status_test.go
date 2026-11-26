package statemachine

import (
	"testing"

	"github.com/Programming-For-the-Community/net-watch/graphql-api/internal/model"
)

func TestStatusValidate_ValidStatusesReturnNil(t *testing.T) {
	valid := []Status{
		{State: model.StateUnknown},
		{State: model.StateUnknown, Failures: defaultThresholds.Failure - 1},
		{State: model.StateUnknown, Successes: defaultThresholds.Recovery - 1},
		{State: model.StateUp, Successes: 100},
		{State: model.StateUp, Failures: defaultThresholds.Failure - 1},
		{State: model.StateDown, Failures: 100},
		{State: model.StateDown, Successes: defaultThresholds.Recovery - 1},
		{State: model.StatePaused},
	}
	for _, s := range valid {
		if err := s.Validate(defaultThresholds); err != nil {
			t.Errorf("%+v.Validate() error = %v, want nil", s, err)
		}
	}
}

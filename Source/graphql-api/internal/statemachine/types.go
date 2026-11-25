package statemachine

import "github.com/Programming-For-the-Community/net-watch/graphql-api/internal/model"

type Thresholds struct {
	Failure, Recovery int
}

type Status struct {
	State     model.CheckState
	Failures  int
	Successes int
}

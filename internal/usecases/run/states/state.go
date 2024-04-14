package states

import (
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --output=./mocks --name=AutomataState
type AutomataState interface {
	Run(ctx context.Context) (StateID, error)
	String() string
}

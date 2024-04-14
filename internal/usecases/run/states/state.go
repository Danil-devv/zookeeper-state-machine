package states

import (
	"context"
	"hw/internal/usecases/run/states/basic"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --output=./mocks --name=AutomataState
type AutomataState interface {
	Run(ctx context.Context) (basic.StateID, error)
	String() string
}

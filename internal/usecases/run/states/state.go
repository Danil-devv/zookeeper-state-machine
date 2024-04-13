package states

import (
	"context"
	"hw/internal/usecases/run/states/basic"
)

type AutomataState interface {
	Run(ctx context.Context) (basic.StateID, error)
	String() string
}

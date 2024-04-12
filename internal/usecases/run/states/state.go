package states

import (
	"context"
	"hw/internal/usecases/run/states/number"
)

type AutomataState interface {
	Run(ctx context.Context) (number.State, error)
	String() string
}

package init

import (
	"context"
	"hw/internal/usecases/run/states/basic"
	"hw/internal/usecases/run/states/number"
	"log/slog"
)

func New(state *basic.State) *State {
	return &State{
		State: state,
	}
}

type State struct {
	*basic.State
}

func (s *State) String() string {
	return "InitState"
}

func (s *State) Run(ctx context.Context) (number.State, error) {
	s.Logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	if ctx.Err() != nil {
		return number.STOPPING, nil
	}
	return number.ATTEMPTER, nil
}

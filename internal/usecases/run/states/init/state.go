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
	if ctx.Err() != nil {
		s.Logger.LogAttrs(
			ctx,
			slog.LevelError,
			"context received an error, stopping",
			slog.String("state", s.String()),
		)

		return number.STOPPING, nil
	}

	s.Logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"switching to the next state",
		slog.String("state", s.String()),
	)
	return number.ATTEMPTER, nil
}

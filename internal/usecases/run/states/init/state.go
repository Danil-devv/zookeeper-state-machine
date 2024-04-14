package init

import (
	"context"
	"hw/internal/usecases/run/states"
	"log/slog"
)

func New(state *states.Basic) *State {
	return &State{
		Basic: state,
	}
}

type State struct {
	*states.Basic
}

func (s *State) String() string {
	return "InitState"
}

func (s *State) Run(ctx context.Context) (states.StateID, error) {
	if ctx.Err() != nil {
		s.Logger.LogAttrs(
			ctx,
			slog.LevelError,
			"context received an error, stopping",
			slog.String("state", s.String()),
		)

		return states.STOPPING, nil
	}
	if !s.Conn.CheckConnection() {
		s.Logger.LogAttrs(
			ctx,
			slog.LevelError,
			"connection failed",
			slog.String("state", s.String()),
		)
		return states.FAILOVER, nil
	}
	return states.ATTEMPTER, nil
}

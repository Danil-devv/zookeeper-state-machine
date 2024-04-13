package init

import (
	"context"
	"hw/internal/usecases/run/states/basic"
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

func (s *State) Run(ctx context.Context) (basic.StateID, error) {
	if ctx.Err() != nil {
		s.Logger.LogAttrs(
			ctx,
			slog.LevelError,
			"context received an error, stopping",
			slog.String("state", s.String()),
		)

		return basic.STOPPING, nil
	}
	if !s.Conn.CheckConnection() {
		s.Logger.LogAttrs(
			ctx,
			slog.LevelError,
			"connection failed",
			slog.String("state", s.String()),
		)
		return basic.FAILOVER, nil
	}
	return basic.ATTEMPTER, nil
}

package stopping

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
	return "StoppingState"
}

func (s *State) Run(ctx context.Context) (basic.StateID, error) {
	s.Logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"closing connection",
		slog.String("state", s.String()),
	)
	s.Conn.Close()
	return basic.EXIT, nil
}

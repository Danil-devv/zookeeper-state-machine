package stopping

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
	return "StoppingState"
}

func (s *State) Run(ctx context.Context) (states.StateID, error) {
	s.Logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"closing connection",
		slog.String("state", s.String()),
	)
	s.Conn.Close()
	return states.EXIT, nil
}

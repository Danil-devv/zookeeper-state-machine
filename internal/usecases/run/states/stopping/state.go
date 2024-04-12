package stopping

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
	return "StoppingState"
}

func (s *State) Run(ctx context.Context) (number.State, error) {
	s.Logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	s.Conn.Close()
	return number.EXIT, nil
}

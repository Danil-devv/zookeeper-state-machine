package stopping

import (
	"context"
	"github.com/go-zookeeper/zk"
	"hw/internal/commands/cmdargs"
	"log/slog"

	"hw/internal/usecases/run/states"
)

func New(log *slog.Logger, args *cmdargs.RunArgs, conn *zk.Conn) *State {
	return &State{
		logger: log,
		args:   args,
		conn:   conn,
	}
}

type State struct {
	logger *slog.Logger
	conn   *zk.Conn
	args   *cmdargs.RunArgs
}

func (s *State) String() string {
	return "StoppingState"
}

func (s *State) Run(ctx context.Context) (states.AutomataState, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	s.conn.Close()
	return nil, nil
}

package init

import (
	"context"
	"github.com/go-zookeeper/zk"
	"hw/internal/commands/cmdargs"
	"hw/internal/usecases/run/states/number"
	"log/slog"
)

func New(logger *slog.Logger, args *cmdargs.RunArgs, conn *zk.Conn) *State {
	return &State{
		logger: logger,
		args:   args,
		conn:   conn,
	}
}

type State struct {
	logger *slog.Logger
	conn   *zk.Conn
	args   *cmdargs.RunArgs
}

func (s *State) GetConn() *zk.Conn {
	return s.conn
}

func (s *State) GetLogger() *slog.Logger {
	return s.logger
}

func (s *State) GetArgs() *cmdargs.RunArgs {
	return s.args
}

func (s *State) String() string {
	return "InitState"
}

func (s *State) Run(ctx context.Context) (number.State, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	if ctx.Err() != nil {
		return number.STOPPING, nil
	}
	return number.ATTEMPTER, nil
}

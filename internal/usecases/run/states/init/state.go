package init

import (
	"context"
	"github.com/go-zookeeper/zk"
	"hw/internal/commands/cmdargs"
	"hw/internal/usecases/run/states"
	"hw/internal/usecases/run/states/attempter"
	"hw/internal/usecases/run/states/stopping"
	"log/slog"
	"time"
)

func New(logger *slog.Logger, args *cmdargs.RunArgs) (*State, error) {
	logger = logger.With("subsystem", "InitState")

	conn, _, err := zk.Connect(args.ZookeeperServers, 3*time.Second)
	if err != nil {
		return nil, err
	}

	return &State{
		logger: logger,
		conn:   conn,
		args:   args,
	}, nil
}

type State struct {
	logger *slog.Logger
	conn   *zk.Conn
	args   *cmdargs.RunArgs
}

func (s *State) String() string {
	return "InitState"
}

func (s *State) Run(ctx context.Context) (states.AutomataState, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	if ctx.Err() != nil {
		return stopping.New(s.logger, s.args, s.conn), nil
	}
	return attempter.New(s.logger, s.args, s.conn), nil
}

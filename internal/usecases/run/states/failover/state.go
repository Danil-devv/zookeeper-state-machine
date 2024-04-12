package failover

import (
	"context"
	"errors"
	"github.com/go-zookeeper/zk"
	"hw/internal/commands/cmdargs"
	"hw/internal/usecases/run/states"
	"hw/internal/usecases/run/states/stopping"
	"log/slog"
	"time"
)

var (
	tryRestartError = errors.New("trying to restart")
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
	return "FailoverState"
}

func (s *State) Run(ctx context.Context) (states.AutomataState, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	ticker := time.NewTicker(s.args.FailoverTimeout)
	for i := 0; i < s.args.FailoverAttemptsCount; i++ {
		select {
		case <-ticker.C:
			_, _, err := zk.Connect(s.args.ZookeeperServers, 3*time.Second)
			if err != nil {
				continue
			}
			return nil, nil
		case <-ctx.Done():
			return stopping.New(s.logger, s.args, s.conn), nil
		}
	}
	return stopping.New(s.logger, s.args, s.conn), nil
}

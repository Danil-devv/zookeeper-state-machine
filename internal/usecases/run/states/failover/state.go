package failover

import (
	"context"
	"github.com/go-zookeeper/zk"
	"hw/internal/commands/cmdargs"
	"hw/internal/usecases/run/states/number"
	"log/slog"
	"time"
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
	return "FailoverState"
}

func (s *State) Run(ctx context.Context) (number.State, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	ticker := time.NewTicker(s.args.FailoverTimeout)
	for i := 0; i < s.args.FailoverAttemptsCount; i++ {
		select {
		case <-ticker.C:
			conn, _, err := zk.Connect(s.args.ZookeeperServers, 3*time.Second)
			if err != nil {
				continue
			}
			s.conn = conn
			return number.INIT, nil
		case <-ctx.Done():
			return number.STOPPING, nil
		}
	}
	return number.STOPPING, nil
}

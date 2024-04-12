package attempter

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
	return "AttempterState"
}

func (s *State) Run(ctx context.Context) (number.State, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	ticker := time.NewTicker(s.args.AttempterTimeout)
	for {
		select {
		case <-ticker.C:
			exists, _, err := s.conn.Exists("/leader")
			if err != nil {
				return number.FAILOVER, nil
			}
			if exists {
				continue
			}
			_, err = s.conn.Create("/leader", []byte("test"), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
			if err != nil {
				return number.FAILOVER, nil
			}
			return number.LEADER, nil
		case <-ctx.Done():
			return number.STOPPING, nil
		}
	}
}

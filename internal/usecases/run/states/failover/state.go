package failover

import (
	"context"
	"github.com/central-university-dev/2024-spring-go-course-lesson8-leader-election/internal/commands/cmdargs"
	"github.com/central-university-dev/2024-spring-go-course-lesson8-leader-election/internal/usecases/run/states"
	"github.com/central-university-dev/2024-spring-go-course-lesson8-leader-election/internal/usecases/run/states/init"
	"github.com/central-university-dev/2024-spring-go-course-lesson8-leader-election/internal/usecases/run/states/stopping"
	"github.com/go-zookeeper/zk"
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

func (s *State) String() string {
	return "FailoverState"
}

func (s *State) Run(ctx context.Context) (states.AutomataState, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	ticker := time.NewTicker(s.args.FailoverTimeout)
	for i := 0; i < s.args.FailoverAttemptsCount; i++ {
		select {
		case <-ticker.C:
			state, err := init.New(s.logger, s.args)
			if err != nil {
				continue
			}
			return state, nil
		case <-ctx.Done():
			return stopping.New(s.logger, s.args, s.conn), nil
		}
	}
	return stopping.New(s.logger, s.args, s.conn), nil
}

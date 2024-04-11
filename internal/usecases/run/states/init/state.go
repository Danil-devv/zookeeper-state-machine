package init

import (
	"context"
	"github.com/central-university-dev/2024-spring-go-course-lesson8-leader-election/internal/commands/cmdargs"
	"github.com/central-university-dev/2024-spring-go-course-lesson8-leader-election/internal/usecases/run/states"
	"github.com/central-university-dev/2024-spring-go-course-lesson8-leader-election/internal/usecases/run/states/attempter"
	"github.com/central-university-dev/2024-spring-go-course-lesson8-leader-election/internal/usecases/run/states/stopping"
	"github.com/go-zookeeper/zk"
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

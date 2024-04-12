package failover

import (
	"context"
	"github.com/go-zookeeper/zk"
	"hw/internal/usecases/run/states/basic"
	"hw/internal/usecases/run/states/number"
	"log/slog"
	"time"
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
	return "FailoverState"
}

func (s *State) Run(ctx context.Context) (number.State, error) {
	s.Logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	ticker := time.NewTicker(s.Args.FailoverTimeout)
	for i := 0; i < s.Args.FailoverAttemptsCount; i++ {
		select {
		case <-ticker.C:
			conn, _, err := zk.Connect(s.Args.ZookeeperServers, 3*time.Second)
			if err != nil {
				continue
			}
			s.Conn = conn
			return number.INIT, nil
		case <-ctx.Done():
			return number.STOPPING, nil
		}
	}
	return number.STOPPING, nil
}

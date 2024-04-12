package attempter

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
	return "AttempterState"
}

func (s *State) Run(ctx context.Context) (number.State, error) {
	s.Logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	ticker := time.NewTicker(s.Args.AttempterTimeout)
	for {
		select {
		case <-ticker.C:
			exists, _, err := s.Conn.Exists("/leader")
			if err != nil {
				return number.FAILOVER, nil
			}
			if exists {
				continue
			}
			_, err = s.Conn.Create("/leader", []byte("test"), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
			if err != nil {
				return number.FAILOVER, nil
			}
			return number.LEADER, nil
		case <-ctx.Done():
			return number.STOPPING, nil
		}
	}
}

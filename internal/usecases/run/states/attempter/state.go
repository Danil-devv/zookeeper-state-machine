package attempter

import (
	"context"
	"hw/internal/usecases/run/states"
	"log/slog"
	"time"

	"github.com/go-zookeeper/zk"
)

func New(state *states.Basic) *State {
	return &State{
		Basic: state,
	}
}

type State struct {
	*states.Basic
}

func (s *State) String() string {
	return "AttempterState"
}

func (s *State) Run(ctx context.Context) (states.StateID, error) {
	ticker := time.NewTicker(s.Args.AttempterTimeout)
	for {
		select {
		case <-ticker.C:
			exists, _, err := s.Conn.Exists("/leader")
			if err != nil {
				s.Logger.LogAttrs(
					ctx,
					slog.LevelError,
					"got an error while trying to check that ephemeral node is exists",
					slog.String("errMsg", err.Error()),
					slog.String("state", s.String()),
				)
				return states.FAILOVER, nil
			}

			if exists {
				s.Logger.LogAttrs(
					ctx,
					slog.LevelInfo,
					"ephemeral node is already exists, continue attempting",
					slog.String("state", s.String()),
				)
				continue
			}

			_, err = s.Conn.Create("/leader", []byte("test"), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
			if err != nil {
				s.Logger.LogAttrs(
					ctx,
					slog.LevelError,
					"cannot create ephemeral node",
					slog.String("errMsg", err.Error()),
					slog.String("state", s.String()),
				)
				return states.FAILOVER, nil
			}
			s.Logger.LogAttrs(
				ctx,
				slog.LevelInfo,
				"successfully create ephemeral node, switching to the next state",
				slog.String("state", s.String()),
			)
			return states.LEADER, nil
		case <-ctx.Done():
			return states.STOPPING, nil
		}
	}
}

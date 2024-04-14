package failover

import (
	"context"
	"hw/internal/usecases/run/states"
	"log/slog"
	"time"
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
	return "FailoverState"
}

func (s *State) Run(ctx context.Context) (states.StateID, error) {
	ticker := time.NewTicker(s.Args.FailoverTimeout)
	for i := 0; i < s.Args.FailoverAttemptsCount; i++ {
		select {
		case <-ticker.C:
			s.Logger.LogAttrs(
				ctx,
				slog.LevelInfo,
				"trying reconnect to zookeeper",
				slog.String("state", s.String()),
			)
			err := s.Conn.Reconnect(s.Args.ZookeeperServers, 3*time.Second)
			if err != nil {
				s.Logger.LogAttrs(
					ctx,
					slog.LevelError,
					"cannot reconnect to zookeeper",
					slog.String("errMsg", err.Error()),
					slog.String("state", s.String()),
				)
				continue
			}
			s.Logger.LogAttrs(
				ctx,
				slog.LevelInfo,
				"successfully reconnected to zookeeper",
				slog.String("state", s.String()),
			)
			return states.INIT, nil
		case <-ctx.Done():
			return states.STOPPING, nil
		}
	}
	return states.STOPPING, nil
}

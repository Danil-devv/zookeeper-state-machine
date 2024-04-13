package failover

import (
	"context"
	"hw/internal/usecases/run/states/basic"
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

func (s *State) Run(ctx context.Context) (basic.StateID, error) {
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
			return basic.INIT, nil
		case <-ctx.Done():
			return basic.STOPPING, nil
		}
	}
	return basic.STOPPING, nil
}

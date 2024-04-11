package leader

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/central-university-dev/2024-spring-go-course-lesson8-leader-election/internal/commands/cmdargs"
	"github.com/central-university-dev/2024-spring-go-course-lesson8-leader-election/internal/usecases/run/states"
	"github.com/central-university-dev/2024-spring-go-course-lesson8-leader-election/internal/usecases/run/states/failover"
	"github.com/central-university-dev/2024-spring-go-course-lesson8-leader-election/internal/usecases/run/states/stopping"
	"github.com/go-zookeeper/zk"
	"github.com/google/uuid"
	"log/slog"
	"path"
	"time"
)

func New(log *slog.Logger, args *cmdargs.RunArgs, conn *zk.Conn) *State {
	return &State{
		logger: log,
		args:   args,
		conn:   conn,
		uuid:   uuid.NewString(),
	}
}

type State struct {
	logger *slog.Logger
	conn   *zk.Conn
	args   *cmdargs.RunArgs
	uuid   string
}

func (s *State) String() string {
	return "LeaderState"
}

func (s *State) Run(ctx context.Context) (states.AutomataState, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	ticker := time.NewTicker(s.args.LeaderTimeout)
	for {
		select {
		case <-ticker.C:
			exists, stat, err := s.conn.Exists(s.args.FileDir)
			if err != nil {
				return failover.New(s.logger, s.args, s.conn), nil
			}

			if exists && int(stat.NumChildren) >= s.args.StorageCapacity {
				childrens, _, err := s.conn.Children(s.args.FileDir)
				if err != nil {
					return failover.New(s.logger, s.args, s.conn), nil
				}
				for i := 0; len(childrens)-i >= s.args.StorageCapacity; i++ {
					err = s.conn.Delete(path.Join(s.args.FileDir, childrens[i]), stat.Version)
					if err != nil {
						return failover.New(s.logger, s.args, s.conn), nil
					}
				}
			}

			if !exists {
				_, err = s.conn.Create(s.args.FileDir, []byte("Leader file directory"), 0, zk.WorldACL(zk.PermAll))
				if err != nil {
					return failover.New(s.logger, s.args, s.conn), nil
				}
			}

			filename, data := s.CreateRandomFile()
			_, err = s.conn.Create(path.Join(s.args.FileDir, filename), data, 0, zk.WorldACL(zk.PermAll))

		case <-ctx.Done():
			return stopping.New(s.logger, s.args, s.conn), nil
		}
	}
}

func (s *State) CreateRandomFile() (name string, data []byte) {
	name = randomdata.Alphanumeric(10)
	data = []byte(fmt.Sprintf("UUID: %s\nText: %s", s.uuid, randomdata.Paragraph()))
	return
}

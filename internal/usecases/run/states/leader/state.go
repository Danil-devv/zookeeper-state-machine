package leader

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/go-zookeeper/zk"
	"github.com/google/uuid"
	"hw/internal/commands/cmdargs"
	"hw/internal/usecases/run/states/number"
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
	return "LeaderState"
}

func (s *State) Run(ctx context.Context) (number.State, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	ticker := time.NewTicker(s.args.LeaderTimeout)
	for {
		select {
		case <-ticker.C:
			exists, stat, err := s.conn.Exists(s.args.FileDir)
			if err != nil {
				return number.FAILOVER, nil
			}

			if exists && int(stat.NumChildren) >= s.args.StorageCapacity {
				childrens, _, err := s.conn.Children(s.args.FileDir)
				if err != nil {
					return number.FAILOVER, nil
				}
				for i := 0; len(childrens)-i >= s.args.StorageCapacity; i++ {
					err = s.conn.Delete(path.Join(s.args.FileDir, childrens[i]), stat.Version)
					if err != nil {
						return number.FAILOVER, nil
					}
				}
			}

			if !exists {
				_, err = s.conn.Create(s.args.FileDir, []byte("Leader file directory"), 0, zk.WorldACL(zk.PermAll))
				if err != nil {
					return number.FAILOVER, nil
				}
			}

			filename, data := s.CreateRandomFile()
			_, err = s.conn.Create(path.Join(s.args.FileDir, filename), data, 0, zk.WorldACL(zk.PermAll))

		case <-ctx.Done():
			return number.STOPPING, nil
		}
	}
}

func (s *State) CreateRandomFile() (name string, data []byte) {
	name = randomdata.Alphanumeric(10)
	data = []byte(fmt.Sprintf("UUID: %s\nText: %s", s.uuid, randomdata.Paragraph()))
	return
}

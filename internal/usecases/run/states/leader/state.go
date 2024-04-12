package leader

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/go-zookeeper/zk"
	"github.com/google/uuid"
	"hw/internal/usecases/run/states/basic"
	"hw/internal/usecases/run/states/number"
	"log/slog"
	"path"
	"time"
)

func New(state *basic.State) *State {
	return &State{
		State: state,
		uuid:  uuid.NewString(),
	}
}

type State struct {
	*basic.State
	uuid string
}

func (s *State) String() string {
	return "LeaderState"
}

func (s *State) Run(ctx context.Context) (number.State, error) {
	s.Logger.LogAttrs(ctx, slog.LevelInfo, "Nothing happened")
	ticker := time.NewTicker(s.Args.LeaderTimeout)
	for {
		select {
		case <-ticker.C:
			exists, stat, err := s.Conn.Exists(s.Args.FileDir)
			if err != nil {
				return number.FAILOVER, nil
			}

			if exists && int(stat.NumChildren) >= s.Args.StorageCapacity {
				childrens, _, err := s.Conn.Children(s.Args.FileDir)
				if err != nil {
					return number.FAILOVER, nil
				}
				for i := 0; len(childrens)-i >= s.Args.StorageCapacity; i++ {
					err = s.Conn.Delete(path.Join(s.Args.FileDir, childrens[i]), stat.Version)
					if err != nil {
						return number.FAILOVER, nil
					}
				}
			}

			if !exists {
				_, err = s.Conn.Create(s.Args.FileDir, []byte("Leader file directory"), 0, zk.WorldACL(zk.PermAll))
				if err != nil {
					return number.FAILOVER, nil
				}
			}

			filename, data := s.CreateRandomFile()
			_, err = s.Conn.Create(path.Join(s.Args.FileDir, filename), data, 0, zk.WorldACL(zk.PermAll))

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

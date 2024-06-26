package leader

import (
	"context"
	"fmt"
	"hw/internal/usecases/run/states"
	"log/slog"
	"path"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/go-zookeeper/zk"
	"github.com/google/uuid"
)

func New(state *states.Basic) *State {
	return &State{
		Basic: state,
		uuid:  uuid.NewString(),
	}
}

type State struct {
	*states.Basic
	uuid string
}

func (s *State) String() string {
	return "LeaderState"
}

func (s *State) Run(ctx context.Context) (states.StateID, error) {
	s.Logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"start writing files in zookeeper",
		slog.String("state", s.String()),
	)
	ticker := time.NewTicker(s.Args.LeaderTimeout)
	for {
		select {
		case <-ticker.C:
			exists, stat, err := s.Conn.Exists(s.Args.FileDir)
			if err != nil {
				s.Logger.LogAttrs(
					ctx,
					slog.LevelError,
					fmt.Sprintf("got an error while trying to check that dir %s is exists", s.Args.FileDir),
					slog.String("errMsg", err.Error()),
					slog.String("state", s.String()),
				)
				return states.FAILOVER, nil
			}

			// TODO: удаляется рандомный ребенок, а не самый старый
			if exists && int(stat.NumChildren) >= s.Args.StorageCapacity {
				s.Logger.LogAttrs(
					ctx,
					slog.LevelInfo,
					fmt.Sprintf("count of childrens exceed maximum, start cleaning dir %s", s.Args.FileDir),
					slog.String("state", s.String()),
				)
				childrens, _, err := s.Conn.Children(s.Args.FileDir)
				if err != nil {
					s.Logger.LogAttrs(
						ctx,
						slog.LevelError,
						"cannot get list of node's childrens",
						slog.String("errMsg", err.Error()),
						slog.String("state", s.String()),
					)
					return states.FAILOVER, nil
				}
				for i := 0; len(childrens)-i >= s.Args.StorageCapacity; i++ {
					err = s.Conn.Delete(path.Join(s.Args.FileDir, childrens[i]), stat.Version)
					if err != nil {
						s.Logger.LogAttrs(
							ctx,
							slog.LevelError,
							"cannot delete child node",
							slog.String("errMsg", err.Error()),
							slog.String("state", s.String()),
						)
						return states.FAILOVER, nil
					}
				}
			}

			if !exists {
				err = s.createFileDir()
				if err != nil {
					s.Logger.LogAttrs(
						ctx,
						slog.LevelError,
						fmt.Sprintf("cannot create working directory %s", s.Args.FileDir),
						slog.String("errMsg", err.Error()),
						slog.String("state", s.String()),
					)
					return states.FAILOVER, nil
				}
			}

			filename, data := s.generateRandomFile()
			s.Logger.LogAttrs(
				ctx,
				slog.LevelInfo,
				"generate random file",
				slog.String("filename", filename),
				slog.String("data", string(data)),
				slog.String("state", s.String()),
			)
			_, err = s.Conn.Create(path.Join(s.Args.FileDir, filename), data, 0, zk.WorldACL(zk.PermAll))
			if err != nil {
				s.Logger.LogAttrs(
					ctx,
					slog.LevelError,
					fmt.Sprintf("cannot create file %s", path.Join(s.Args.FileDir, filename)),
					slog.String("data", string(data)),
					slog.String("errMsg", err.Error()),
					slog.String("state", s.String()),
				)
				return states.FAILOVER, nil
			}

		case <-ctx.Done():
			return states.STOPPING, nil
		}
	}
}

func (s *State) generateRandomFile() (name string, data []byte) {
	name = randomdata.Alphanumeric(10)
	data = []byte(fmt.Sprintf("UUID: %s\nText: %s", s.uuid, randomdata.Paragraph()))
	return
}

func (s *State) createFileDir() (err error) {
	nodes := strings.Split(strings.Trim(s.Args.FileDir, "/"), "/")
	curPath := ""
	f := false
	for _, node := range nodes {
		curPath += "/" + node
		if f {
			_, err = s.Conn.Create(curPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
			if err != nil {
				return err
			}
			continue
		}

		exists, _, err := s.Conn.Exists(curPath)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		f = true
		_, err = s.Conn.Create(curPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}
	return nil
}

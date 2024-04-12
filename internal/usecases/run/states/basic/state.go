package basic

import (
	"github.com/go-zookeeper/zk"
	"hw/internal/commands/cmdargs"
	"log/slog"
)

func New(log *slog.Logger, args *cmdargs.RunArgs, conn *zk.Conn) *State {
	return &State{
		Logger: log,
		Conn:   conn,
		Args:   args,
	}
}

type State struct {
	Logger *slog.Logger
	Conn   *zk.Conn
	Args   *cmdargs.RunArgs
}

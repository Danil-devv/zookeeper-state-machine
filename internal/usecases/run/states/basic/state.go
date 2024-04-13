package basic

import (
	"hw/internal/commands/cmdargs"
	"log/slog"
)

func New(log *slog.Logger, args *cmdargs.RunArgs, conn ZkConn) *State {
	return &State{
		Logger: log,
		Conn:   conn,
		Args:   args,
	}
}

type State struct {
	Logger *slog.Logger
	Conn   ZkConn
	Args   *cmdargs.RunArgs
}

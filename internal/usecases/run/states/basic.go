package states

import (
	"hw/internal/commands/cmdargs"
	"log/slog"
)

func NewBasicState(log *slog.Logger, args *cmdargs.RunArgs, conn ZkConn) *Basic {
	return &Basic{
		Logger: log,
		Conn:   conn,
		Args:   args,
	}
}

type Basic struct {
	Logger *slog.Logger
	Conn   ZkConn
	Args   *cmdargs.RunArgs
}

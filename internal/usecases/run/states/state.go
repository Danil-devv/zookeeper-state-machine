package states

import (
	"context"
	"github.com/go-zookeeper/zk"
	"hw/internal/commands/cmdargs"
	"hw/internal/usecases/run/states/number"
	"log/slog"
)

type AutomataState interface {
	Run(ctx context.Context) (number.State, error)
	GetConn() *zk.Conn
	GetLogger() *slog.Logger
	GetArgs() *cmdargs.RunArgs
	String() string
}

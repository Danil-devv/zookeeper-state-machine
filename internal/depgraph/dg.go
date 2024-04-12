package depgraph

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"hw/internal/commands/cmdargs"
	"hw/internal/usecases/run"
	"hw/internal/usecases/run/states/basic"
	initstate "hw/internal/usecases/run/states/init"
	"log/slog"
	"os"
	"sync"
	"time"
)

type dgEntity[T any] struct {
	sync.Once
	value   T
	initErr error
}

func (e *dgEntity[T]) get(init func() (T, error)) (T, error) {
	e.Do(func() {
		e.value, e.initErr = init()
	})
	if e.initErr != nil {
		return *new(T), e.initErr
	}
	return e.value, nil
}

type DepGraph struct {
	logger      *dgEntity[*slog.Logger]
	stateRunner *dgEntity[*run.LoopRunner]
	initState   *dgEntity[*initstate.State]
	zkConn      *dgEntity[*zk.Conn]
	basicState  *dgEntity[*basic.State]
}

func New() *DepGraph {
	return &DepGraph{
		logger:      &dgEntity[*slog.Logger]{},
		stateRunner: &dgEntity[*run.LoopRunner]{},
		initState:   &dgEntity[*initstate.State]{},
		zkConn:      &dgEntity[*zk.Conn]{},
		basicState:  &dgEntity[*basic.State]{},
	}
}

func (dg *DepGraph) GetLogger() (*slog.Logger, error) {
	return dg.logger.get(func() (*slog.Logger, error) {
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})), nil
	})
}

func (dg *DepGraph) GetRunner() (run.Runner, error) {
	return dg.stateRunner.get(func() (*run.LoopRunner, error) {
		logger, err := dg.GetLogger()
		if err != nil {
			return nil, fmt.Errorf("get logger: %w", err)
		}
		return run.NewLoopRunner(logger), nil
	})
}

func (dg *DepGraph) GetZkConn(args *cmdargs.RunArgs) (*zk.Conn, error) {
	return dg.zkConn.get(func() (*zk.Conn, error) {
		conn, _, err := zk.Connect(args.ZookeeperServers, 3*time.Second)
		if err != nil {
			return nil, err
		}
		return conn, nil
	})
}

func (dg *DepGraph) GetInitState(state *basic.State) (*initstate.State, error) {
	return dg.initState.get(func() (*initstate.State, error) {
		return initstate.New(state), nil
	})
}

func (dg *DepGraph) GetBasicState(conn *zk.Conn, args *cmdargs.RunArgs, l *slog.Logger) (*basic.State, error) {
	return dg.basicState.get(func() (*basic.State, error) {
		return basic.New(l, args, conn), nil
	})
}

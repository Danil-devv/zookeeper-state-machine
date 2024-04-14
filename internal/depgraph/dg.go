package depgraph

import (
	"fmt"
	"hw/internal/adapters/zookeeper"
	"hw/internal/commands/cmdargs"
	"hw/internal/usecases/run"
	"hw/internal/usecases/run/states"
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
	zkConn      *dgEntity[*zookeeper.Conn]
	basicState  *dgEntity[*states.Basic]
}

func New() *DepGraph {
	return &DepGraph{
		logger:      &dgEntity[*slog.Logger]{},
		stateRunner: &dgEntity[*run.LoopRunner]{},
		initState:   &dgEntity[*initstate.State]{},
		zkConn:      &dgEntity[*zookeeper.Conn]{},
		basicState:  &dgEntity[*states.Basic]{},
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

func (dg *DepGraph) GetZkConn(args *cmdargs.RunArgs) (*zookeeper.Conn, error) {
	return dg.zkConn.get(func() (*zookeeper.Conn, error) {
		conn, err := zookeeper.NewConn(args.ZookeeperServers, 3*time.Second)
		if err != nil {
			return nil, err
		}
		return conn, nil
	})
}

func (dg *DepGraph) GetInitState(state *states.Basic) (*initstate.State, error) {
	return dg.initState.get(func() (*initstate.State, error) {
		return initstate.New(state), nil
	})
}

func (dg *DepGraph) GetBasicState(args *cmdargs.RunArgs, l *slog.Logger, conn *zookeeper.Conn) (*states.Basic, error) {
	return dg.basicState.get(func() (*states.Basic, error) {
		return states.NewBasicState(l, args, conn), nil
	})
}

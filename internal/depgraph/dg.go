package depgraph

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"hw/internal/commands/cmdargs"
	"hw/internal/usecases/run"
	initialization "hw/internal/usecases/run/states/init"
	"log/slog"
	"os"
	"sync"
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
	initState   *dgEntity[*initialization.State]
	zkConn      *dgEntity[*zk.Conn]
}

func New() *DepGraph {
	return &DepGraph{
		logger:      &dgEntity[*slog.Logger]{},
		stateRunner: &dgEntity[*run.LoopRunner]{},
		initState:   &dgEntity[*initialization.State]{},
		zkConn:      &dgEntity[*zk.Conn]{},
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

func (dg *DepGraph) GetInitState(args *cmdargs.RunArgs) (*initialization.State, error) {
	return dg.initState.get(func() (*initialization.State, error) {
		logger, err := dg.GetLogger()
		if err != nil {
			return nil, fmt.Errorf("get logger: %w", err)
		}
		state, err := initialization.New(logger, args)
		if err != nil {
			return nil, err
		}
		return state, nil
	})
}

package run

import (
	"context"
	"fmt"
	"hw/internal/usecases/run/states/attempter"
	"hw/internal/usecases/run/states/failover"
	initstate "hw/internal/usecases/run/states/init"
	"hw/internal/usecases/run/states/leader"
	"hw/internal/usecases/run/states/number"
	"hw/internal/usecases/run/states/stopping"
	"log/slog"

	"hw/internal/usecases/run/states"
)

var _ Runner = &LoopRunner{}

type Runner interface {
	Run(ctx context.Context, state states.AutomataState) error
}

func NewLoopRunner(logger *slog.Logger) *LoopRunner {
	logger = logger.With("subsystem", "StateRunner")
	return &LoopRunner{
		logger: logger,
	}
}

type LoopRunner struct {
	logger *slog.Logger
}

func (r *LoopRunner) Run(ctx context.Context, state states.AutomataState) error {
	for state != nil {
		r.logger.LogAttrs(ctx, slog.LevelInfo, "start running state", slog.String("state", state.String()))
		var err error
		state, err = processState(ctx, state)
		if err != nil {
			return fmt.Errorf("state %s run: %w", state.String(), err)
		}
	}
	r.logger.LogAttrs(ctx, slog.LevelInfo, "no new state, finish")
	return nil
}

func processState(ctx context.Context, state states.AutomataState) (states.AutomataState, error) {
	n, err := state.Run(ctx)
	if err != nil {
		return nil, err
	}

	switch n {
	case number.INIT:
		state = initstate.New(state.GetLogger(), state.GetArgs(), state.GetConn())
	case number.STOPPING:
		state = stopping.New(state.GetLogger(), state.GetArgs(), state.GetConn())
	case number.FAILOVER:
		state = failover.New(state.GetLogger(), state.GetArgs(), state.GetConn())
	case number.ATTEMPTER:
		state = attempter.New(state.GetLogger(), state.GetArgs(), state.GetConn())
	case number.LEADER:
		state = leader.New(state.GetLogger(), state.GetArgs(), state.GetConn())
	case number.EXIT:
		state = nil
	}

	return state, nil
}

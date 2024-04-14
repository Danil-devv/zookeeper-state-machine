package run

import (
	"context"
	"fmt"
	"hw/internal/usecases/run/states/attempter"
	"hw/internal/usecases/run/states/failover"
	initstate "hw/internal/usecases/run/states/init"
	"hw/internal/usecases/run/states/leader"
	"hw/internal/usecases/run/states/stopping"
	"log/slog"

	"hw/internal/usecases/run/states"
)

var _ Runner = &LoopRunner{}

type Runner interface {
	Run(ctx context.Context, state states.AutomataState, b *states.Basic) error
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

func (r *LoopRunner) Run(ctx context.Context, state states.AutomataState, b *states.Basic) error {
	for state != nil {
		r.logger.LogAttrs(
			ctx,
			slog.LevelInfo,
			"start running state",
			slog.String("state", state.String()),
		)

		var err error
		state, err = processState(ctx, state, b)
		if err != nil {
			return fmt.Errorf("state %s run: %w", state.String(), err)
		}
	}
	r.logger.LogAttrs(ctx, slog.LevelInfo, "no new state, finish")
	return nil
}

func processState(ctx context.Context, state states.AutomataState, b *states.Basic) (states.AutomataState, error) {
	n, err := state.Run(ctx)
	if err != nil {
		return nil, err
	}

	switch n {
	case states.INIT:
		state = initstate.New(b)
	case states.STOPPING:
		state = stopping.New(b)
	case states.FAILOVER:
		state = failover.New(b)
	case states.ATTEMPTER:
		state = attempter.New(b)
	case states.LEADER:
		state = leader.New(b)
	case states.EXIT:
		state = nil
	}

	return state, nil
}

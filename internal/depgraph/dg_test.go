package depgraph

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"hw/internal/adapters/zookeeper"
	"hw/internal/commands/cmdargs"
	"hw/internal/usecases/run"
	"hw/internal/usecases/run/states/basic"
	initstate "hw/internal/usecases/run/states/init"
	"log/slog"
	"testing"
)

func TestNew(t *testing.T) {
	dg := New()
	assert.Equal(t, dg.logger, &dgEntity[*slog.Logger]{})
	assert.Equal(t, dg.stateRunner, &dgEntity[*run.LoopRunner]{})
	assert.Equal(t, dg.initState, &dgEntity[*initstate.State]{})
	assert.Equal(t, dg.basicState, &dgEntity[*basic.State]{})
	assert.Equal(t, dg.zkConn, &dgEntity[*zookeeper.Conn]{})

	assert.NoError(t, dg.logger.initErr)
	assert.NoError(t, dg.stateRunner.initErr)
	assert.NoError(t, dg.initState.initErr)
	assert.NoError(t, dg.basicState.initErr)
	assert.NoError(t, dg.zkConn.initErr)
}

func TestEntity_MethodGetOnce(t *testing.T) {
	type testcase struct {
		v int
		e error
	}
	tests := []testcase{
		{8, nil},
		{10, nil},
		{18, errors.New("some error")},
		{-5, nil},
		{0, errors.New("some error")},
	}

	for _, tc := range tests {
		e := &dgEntity[int]{}
		v, err := e.get(func() (int, error) {
			return tc.v, tc.e
		})
		assert.Equal(t, tc.e, err)
		if err != nil {
			assert.Equal(t, 0, v)
		} else {
			assert.Equal(t, tc.v, v)
		}
	}
}

func TestEntity_MethodGetMultiply(t *testing.T) {
	type testcase struct {
		v int
		e error
	}
	tests := []testcase{
		{8, nil},
		{10, nil},
		{18, errors.New("some error")},
		{-5, nil},
		{0, errors.New("some error")},
	}

	e := &dgEntity[int]{}
	for _, tc := range tests {
		v, err := e.get(func() (int, error) {
			return tc.v, tc.e
		})
		assert.Equal(t, tests[0].v, v)
		assert.Equal(t, tests[0].e, err)
	}
}

func TestDepGraph_GetCustomLogger_NoError(t *testing.T) {
	myLogger := &slog.Logger{}
	customLogger := &dgEntity[*slog.Logger]{}
	_, _ = customLogger.get(func() (*slog.Logger, error) {
		return myLogger, nil
	})

	dg := New()
	dg.logger = customLogger

	logger1, err := dg.GetLogger()
	assert.NoError(t, err)
	assert.Equal(t, myLogger, logger1)

	logger2, err := dg.GetLogger()
	assert.NoError(t, err)
	assert.Equal(t, myLogger, logger2)
	assert.Equal(t, logger1, logger2)
}

func TestDepGraph_GetCustomLogger_WithError(t *testing.T) {
	dg := New()

	customLogger := &dgEntity[*slog.Logger]{}
	_, _ = customLogger.get(func() (*slog.Logger, error) {
		return nil, errors.New("some error")
	})
	dg.logger = customLogger
	logger, err := dg.GetLogger()
	assert.Error(t, err)
	assert.Empty(t, logger)

	logger, err = dg.GetLogger()
	assert.Error(t, err)
	assert.Empty(t, logger)

	customLogger = &dgEntity[*slog.Logger]{}
	_, _ = customLogger.get(func() (*slog.Logger, error) {
		return &slog.Logger{}, errors.New("some error")
	})
	dg.logger = customLogger
	logger, err = dg.GetLogger()
	assert.Error(t, err)
	assert.Empty(t, logger)

	logger, err = dg.GetLogger()
	assert.Error(t, err)
	assert.Empty(t, logger)
}

func TestDepGraph_GetLogger(t *testing.T) {
	dg := New()

	logger, err := dg.GetLogger()
	assert.NoError(t, err)
	assert.NotEmpty(t, logger)
}

func TestDepGraph_GetCustomRunner_NoError(t *testing.T) {
	myRunner := &run.LoopRunner{}
	customRunner := &dgEntity[*run.LoopRunner]{}
	_, _ = customRunner.get(func() (*run.LoopRunner, error) {
		return myRunner, nil
	})

	dg := New()
	dg.stateRunner = customRunner

	runner1, err := dg.GetRunner()
	assert.NoError(t, err)
	assert.Equal(t, myRunner, runner1)

	runner2, err := dg.GetRunner()
	assert.NoError(t, err)
	assert.Equal(t, myRunner, runner2)
	assert.Equal(t, runner1, runner2)
}

func TestDepGraph_GetCustomRunner_WithError(t *testing.T) {
	myRunner := &run.LoopRunner{}
	customRunner := &dgEntity[*run.LoopRunner]{}
	_, _ = customRunner.get(func() (*run.LoopRunner, error) {
		return myRunner, errors.New("some error")
	})

	dg := New()
	dg.stateRunner = customRunner
	runner, err := dg.GetRunner()
	assert.Error(t, err)
	assert.Empty(t, runner)

	dg = New()
	customLogger := &dgEntity[*slog.Logger]{}
	_, _ = customLogger.get(func() (*slog.Logger, error) {
		return nil, errors.New("some error")
	})
	dg.logger = customLogger
	runner, err = dg.GetRunner()
	assert.Error(t, err)
	assert.Empty(t, runner)
}

func TestDepGraph_GetRunner(t *testing.T) {
	dg := New()

	runner, err := dg.GetRunner()
	assert.NoError(t, err)
	assert.NotEmpty(t, runner)
}

func TestDepGraph_GetZkConn_Wrong(t *testing.T) {
	dg := New()
	args := &cmdargs.RunArgs{ZookeeperServers: []string{"wronghost:wrongport"}}
	conn, err := dg.GetZkConn(args)
	assert.Error(t, err)
	assert.Empty(t, conn)
}

func TestDepGraph_GetInitState(t *testing.T) {
	dg := New()
	b := &basic.State{}
	state, err := dg.GetInitState(b)
	assert.NoError(t, err)
	assert.NotEmpty(t, state)
}

func TestDepGraph_GetBasicState(t *testing.T) {
	dg := New()
	args := &cmdargs.RunArgs{}
	logger := &slog.Logger{}
	conn := &zookeeper.Conn{}
	state, err := dg.GetBasicState(args, logger, conn)
	assert.NoError(t, err)
	assert.Equal(t, args, state.Args)
	assert.Equal(t, logger, state.Logger)
	assert.Equal(t, conn, state.Conn)
}

package run

import (
	"context"
	"errors"
	"hw/internal/usecases/run/states"
	"hw/internal/usecases/run/states/attempter"
	"hw/internal/usecases/run/states/failover"
	initstate "hw/internal/usecases/run/states/init"
	"hw/internal/usecases/run/states/mocks"
	"hw/internal/usecases/run/states/stopping"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessState(t *testing.T) {
	basic := &states.Basic{}
	ctx := context.Background()

	attempterState := mocks.NewAutomataState(t)
	attempterState.On("Run", ctx).Return(states.ATTEMPTER, nil)

	initState := mocks.NewAutomataState(t)
	initState.On("Run", ctx).Return(states.INIT, nil)

	failoverState := mocks.NewAutomataState(t)
	failoverState.On("Run", ctx).Return(states.FAILOVER, nil)

	stoppingState := mocks.NewAutomataState(t)
	stoppingState.On("Run", ctx).Return(states.STOPPING, nil)

	exit := mocks.NewAutomataState(t)
	exit.On("Run", ctx).Return(states.EXIT, nil)

	stateWithErr := mocks.NewAutomataState(t)
	stateWithErr.On("Run", ctx).Return(states.ATTEMPTER, errors.New("some error"))
	tests := []struct {
		name          string
		state         states.AutomataState
		expectedState states.AutomataState
		expectedError error
	}{
		{
			name:  "Attempter state",
			state: attempterState,
			expectedState: &attempter.State{
				Basic: basic,
			},
			expectedError: nil,
		},
		{
			name:  "Init state",
			state: initState,
			expectedState: &initstate.State{
				Basic: basic,
			},
			expectedError: nil,
		},
		{
			name:  "Failover state",
			state: failoverState,
			expectedState: &failover.State{
				Basic: basic,
			},
			expectedError: nil,
		},
		{
			name:  "Stopping state",
			state: stoppingState,
			expectedState: &stopping.State{
				Basic: basic,
			},
			expectedError: nil,
		},
		{
			name:          "Exit state",
			state:         exit,
			expectedState: nil,
			expectedError: nil,
		},
		{
			name:          "State with error",
			state:         stateWithErr,
			expectedState: nil,
			expectedError: errors.New("some error"),
		},
	}

	for _, tc := range tests {
		state, err := processState(ctx, tc.state, basic)
		assert.Equal(t, tc.expectedState, state)
		assert.Equal(t, tc.expectedError, err)
	}
}

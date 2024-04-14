// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"
	basic "hw/internal/usecases/run/states/basic"

	mock "github.com/stretchr/testify/mock"
)

// AutomataState is an autogenerated mock type for the AutomataState type
type AutomataState struct {
	mock.Mock
}

// Run provides a mock function with given fields: ctx
func (_m *AutomataState) Run(ctx context.Context) (basic.StateID, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Run")
	}

	var r0 basic.StateID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (basic.StateID, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) basic.StateID); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(basic.StateID)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// String provides a mock function with given fields:
func (_m *AutomataState) String() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for String")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewAutomataState creates a new instance of AutomataState. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAutomataState(t interface {
	mock.TestingT
	Cleanup(func())
}) *AutomataState {
	mock := &AutomataState{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

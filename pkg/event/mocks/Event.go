// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	event "github.com/veleton777/booking_api/pkg/event"
)

// Event is an autogenerated mock type for the Event type
type Event struct {
	mock.Mock
}

// SaveEvent provides a mock function with given fields: ctx, entity
func (_m *Event) SaveEvent(ctx context.Context, entity event.Entity) error {
	ret := _m.Called(ctx, entity)

	if len(ret) == 0 {
		panic("no return value specified for SaveEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, event.Entity) error); ok {
		r0 = rf(ctx, entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewEvent creates a new instance of Event. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEvent(t interface {
	mock.TestingT
	Cleanup(func())
}) *Event {
	mock := &Event{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

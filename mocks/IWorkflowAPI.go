// Code generated by mockery v1.1.1. DO NOT EDIT.

package mocks

import (
	codefresh "github.com/codefresh-io/go-sdk/pkg/codefresh"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// IWorkflowAPI is an autogenerated mock type for the IWorkflowAPI type
type IWorkflowAPI struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0
func (_m *IWorkflowAPI) Get(_a0 string) (*codefresh.Workflow, error) {
	ret := _m.Called(_a0)

	var r0 *codefresh.Workflow
	if rf, ok := ret.Get(0).(func(string) *codefresh.Workflow); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*codefresh.Workflow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WaitForStatus provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *IWorkflowAPI) WaitForStatus(_a0 string, _a1 string, _a2 time.Duration, _a3 time.Duration) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, time.Duration, time.Duration) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Code generated by mockery v1.1.1. DO NOT EDIT.

package mocks

import (
	model "github.com/codefresh-io/argo-platform/libs/ql/graph/model"
	mock "github.com/stretchr/testify/mock"
)

// IArgoRuntimeAPI is an autogenerated mock type for the IArgoRuntimeAPI type
type IArgoRuntimeAPI struct {
	mock.Mock
}

// List provides a mock function with given fields:
func (_m *IArgoRuntimeAPI) List() ([]model.Runtime, error) {
	ret := _m.Called()

	var r0 []model.Runtime
	if rf, ok := ret.Get(0).(func() []model.Runtime); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Runtime)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	logrus "github.com/sirupsen/logrus"
	mock "github.com/stretchr/testify/mock"

	runner "github.com/weiliang-ms/easyctl/pkg/runner"
)

// HandleOSInterface is an autogenerated mock type for the HandleOSInterface type
type HandleOSInterface struct {
	mock.Mock
}

// DoRequest provides a mock function with given fields: doRequestItem
func (_m *HandleOSInterface) DoRequest(doRequestItem runner.DoRequestItem) (string, error) {
	ret := _m.Called(doRequestItem)

	var r0 string
	if rf, ok := ret.Get(0).(func(runner.DoRequestItem) string); ok {
		r0 = rf(doRequestItem)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(runner.DoRequestItem) error); ok {
		r1 = rf(doRequestItem)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCPUInfo provides a mock function with given fields: s, l
func (_m *HandleOSInterface) GetCPUInfo(s runner.ServerInternal, l *logrus.Logger) (string, error) {
	ret := _m.Called(s, l)

	var r0 string
	if rf, ok := ret.Get(0).(func(runner.ServerInternal, *logrus.Logger) string); ok {
		r0 = rf(s, l)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(runner.ServerInternal, *logrus.Logger) error); ok {
		r1 = rf(s, l)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCPULoadAverage provides a mock function with given fields: s, l
func (_m *HandleOSInterface) GetCPULoadAverage(s runner.ServerInternal, l *logrus.Logger) (string, error) {
	ret := _m.Called(s, l)

	var r0 string
	if rf, ok := ret.Get(0).(func(runner.ServerInternal, *logrus.Logger) string); ok {
		r0 = rf(s, l)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(runner.ServerInternal, *logrus.Logger) error); ok {
		r1 = rf(s, l)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHostName provides a mock function with given fields: s, l
func (_m *HandleOSInterface) GetHostName(s runner.ServerInternal, l *logrus.Logger) (string, error) {
	ret := _m.Called(s, l)

	var r0 string
	if rf, ok := ret.Get(0).(func(runner.ServerInternal, *logrus.Logger) string); ok {
		r0 = rf(s, l)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(runner.ServerInternal, *logrus.Logger) error); ok {
		r1 = rf(s, l)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetKernelVersion provides a mock function with given fields: s, l
func (_m *HandleOSInterface) GetKernelVersion(s runner.ServerInternal, l *logrus.Logger) (string, error) {
	ret := _m.Called(s, l)

	var r0 string
	if rf, ok := ret.Get(0).(func(runner.ServerInternal, *logrus.Logger) string); ok {
		r0 = rf(s, l)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(runner.ServerInternal, *logrus.Logger) error); ok {
		r1 = rf(s, l)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMemoryInfo provides a mock function with given fields: s, l
func (_m *HandleOSInterface) GetMemoryInfo(s runner.ServerInternal, l *logrus.Logger) (string, error) {
	ret := _m.Called(s, l)

	var r0 string
	if rf, ok := ret.Get(0).(func(runner.ServerInternal, *logrus.Logger) string); ok {
		r0 = rf(s, l)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(runner.ServerInternal, *logrus.Logger) error); ok {
		r1 = rf(s, l)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMountPointInfo provides a mock function with given fields: s, l
func (_m *HandleOSInterface) GetMountPointInfo(s runner.ServerInternal, l *logrus.Logger) (string, error) {
	ret := _m.Called(s, l)

	var r0 string
	if rf, ok := ret.Get(0).(func(runner.ServerInternal, *logrus.Logger) string); ok {
		r0 = rf(s, l)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(runner.ServerInternal, *logrus.Logger) error); ok {
		r1 = rf(s, l)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSystemVersion provides a mock function with given fields: s, l
func (_m *HandleOSInterface) GetSystemVersion(s runner.ServerInternal, l *logrus.Logger) (string, error) {
	ret := _m.Called(s, l)

	var r0 string
	if rf, ok := ret.Get(0).(func(runner.ServerInternal, *logrus.Logger) string); ok {
		r0 = rf(s, l)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(runner.ServerInternal, *logrus.Logger) error); ok {
		r1 = rf(s, l)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

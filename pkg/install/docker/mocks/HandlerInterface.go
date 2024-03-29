// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	logrus "github.com/sirupsen/logrus"
	install "github.com/weiliang-ms/easyctl/pkg/install"

	mock "github.com/stretchr/testify/mock"

	runner "github.com/weiliang-ms/easyctl/pkg/runner"

	time "time"
)

// HandlerInterface is an autogenerated mock type for the HandlerInterface type
type HandlerInterface struct {
	mock.Mock
}

// Boot provides a mock function with given fields: server, local, logger, timeout
func (_m *HandlerInterface) Boot(server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.BootErr {
	ret := _m.Called(server, local, logger, timeout)

	var r0 install.BootErr
	if rf, ok := ret.Get(0).(func(runner.ServerInternal, bool, *logrus.Logger, time.Duration) install.BootErr); ok {
		r0 = rf(server, local, logger, timeout)
	} else {
		r0 = ret.Get(0).(install.BootErr)
	}

	return r0
}

// Detect provides a mock function with given fields: cmd, server, local, logger, timeout
func (_m *HandlerInterface) Detect(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) error {
	ret := _m.Called(cmd, server, local, logger, timeout)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, runner.ServerInternal, bool, *logrus.Logger, time.Duration) error); ok {
		r0 = rf(cmd, server, local, logger, timeout)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exec provides a mock function with given fields: cmd, server, local, logger, timeout
func (_m *HandlerInterface) Exec(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) error {
	ret := _m.Called(cmd, server, local, logger, timeout)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, runner.ServerInternal, bool, *logrus.Logger, time.Duration) error); ok {
		r0 = rf(cmd, server, local, logger, timeout)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HandPackage provides a mock function with given fields: server, filePath, local, logger, timeout
func (_m *HandlerInterface) HandPackage(server runner.ServerInternal, filePath string, local bool, logger *logrus.Logger, timeout time.Duration) error {
	ret := _m.Called(server, filePath, local, logger, timeout)

	var r0 error
	if rf, ok := ret.Get(0).(func(runner.ServerInternal, string, bool, *logrus.Logger, time.Duration) error); ok {
		r0 = rf(server, filePath, local, logger, timeout)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Install provides a mock function with given fields: cmd, server, local, logger, timeout
func (_m *HandlerInterface) Install(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.InstallErr {
	ret := _m.Called(cmd, server, local, logger, timeout)

	var r0 install.InstallErr
	if rf, ok := ret.Get(0).(func(string, runner.ServerInternal, bool, *logrus.Logger, time.Duration) install.InstallErr); ok {
		r0 = rf(cmd, server, local, logger, timeout)
	} else {
		r0 = ret.Get(0).(install.InstallErr)
	}

	return r0
}

// Prune provides a mock function with given fields: server, local, logger, timeout
func (_m *HandlerInterface) Prune(server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.PruneErr {
	ret := _m.Called(server, local, logger, timeout)

	var r0 install.PruneErr
	if rf, ok := ret.Get(0).(func(runner.ServerInternal, bool, *logrus.Logger, time.Duration) install.PruneErr); ok {
		r0 = rf(server, local, logger, timeout)
	} else {
		r0 = ret.Get(0).(install.PruneErr)
	}

	return r0
}

// SetConfig provides a mock function with given fields: cmd, server, local, logger, timeout
func (_m *HandlerInterface) SetConfig(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.SetConfigErr {
	ret := _m.Called(cmd, server, local, logger, timeout)

	var r0 install.SetConfigErr
	if rf, ok := ret.Get(0).(func(string, runner.ServerInternal, bool, *logrus.Logger, time.Duration) install.SetConfigErr); ok {
		r0 = rf(cmd, server, local, logger, timeout)
	} else {
		r0 = ret.Get(0).(install.SetConfigErr)
	}

	return r0
}

// SetSystemd provides a mock function with given fields: cmd, server, local, logger, timeout
func (_m *HandlerInterface) SetSystemd(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.SetSystemdErr {
	ret := _m.Called(cmd, server, local, logger, timeout)

	var r0 install.SetSystemdErr
	if rf, ok := ret.Get(0).(func(string, runner.ServerInternal, bool, *logrus.Logger, time.Duration) install.SetSystemdErr); ok {
		r0 = rf(cmd, server, local, logger, timeout)
	} else {
		r0 = ret.Get(0).(install.SetSystemdErr)
	}

	return r0
}

// SetUpRuntime provides a mock function with given fields: cmd, server, local, logger, timeout
func (_m *HandlerInterface) SetUpRuntime(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.SetUpRuntimeErr {
	ret := _m.Called(cmd, server, local, logger, timeout)

	var r0 install.SetUpRuntimeErr
	if rf, ok := ret.Get(0).(func(string, runner.ServerInternal, bool, *logrus.Logger, time.Duration) install.SetUpRuntimeErr); ok {
		r0 = rf(cmd, server, local, logger, timeout)
	} else {
		r0 = ret.Get(0).(install.SetUpRuntimeErr)
	}

	return r0
}

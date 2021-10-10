package set

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"testing"
)

func TestNewPasswordScript(t *testing.T) {
	// test parse error
	b := `newRootPassword: 
			- "Cloud@2021%^&*"`
	content, err := NewPasswordScript([]byte(b), NewPasswordTmpl)
	assert.Equal(t, "", content)
	assert.EqualError(t, err, "yaml: line 2: found character that cannot start any token")

	// test set weak password
	b = `newRootPassword: 123`
	content, err = NewPasswordScript([]byte(b), NewPasswordTmpl)
	assert.EqualError(t, err, "密码长度：3 不符合标准")

	// test valid config
	b = `newRootPassword: 123456`
	content, err = NewPasswordScript([]byte(b), NewPasswordTmpl)
	assert.Nil(t, err)
	assert.Equal(t, "\n#!/bin/bash\nset -e\necho \"123456\" | passwd --stdin root\n", content)
}
func TestNewPassword(t *testing.T) {
	b := `newRootPassword: 
			- "Cloud@2021%^&*"`
	item := command.OperationItem{
		B:      []byte(b),
		Logger: logrus.New(),
	}
	err := NewPassword(item)
	assert.EqualError(t, err, "yaml: line 2: found character that cannot start any token")

	item.B = []byte(`newRootPassword: "Cloud@2021%^&*"`)
	assert.Equal(t, command.RunErr{}, NewPassword(item))
}

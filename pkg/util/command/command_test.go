package command

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetExecutorDefault(t *testing.T) {
	// test -c is nil
	entity := ExecutorEntity{}
	assert.Nil(t, SetExecutorDefault(entity, ""))

	// test debug flag
	entity = ExecutorEntity{}
	entity.Cmd = &cobra.Command{}
	parentCmd := &cobra.Command{}
	parentCmd.AddCommand(entity.Cmd)
	parentParentCmd := &cobra.Command{}
	parentParentCmd.AddCommand(parentCmd)
	assert.EqualError(t, SetExecutorDefault(entity, "config.yaml"), "flag accessed but not defined: debug")

	// test err path configfile & full in DefaultConfig
	var debug bool

	entity = ExecutorEntity{}
	entity.Cmd = &cobra.Command{}
	entity.DefaultConfig = []byte("ddd")
	parentCmd = &cobra.Command{}
	parentCmd.AddCommand(entity.Cmd)
	parentParentCmd = &cobra.Command{}
	parentParentCmd.AddCommand(parentCmd)
	parentParentCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug")

	assert.Errorf(t, SetExecutorDefault(entity, "1.txt"), "open 1.txt: The system cannot find the file specified.")

}

// test panic return
func TestSetExecutorDefaultReturnErr(t *testing.T) {
	var debug bool
	entity := ExecutorEntity{}
	entity.Cmd = &cobra.Command{}
	parentCmd := &cobra.Command{}
	parentCmd.AddCommand(entity.Cmd)
	parentParentCmd := &cobra.Command{}
	parentParentCmd.AddCommand(parentCmd)
	parentParentCmd.PersistentFlags().BoolVar(&debug, "debug", true, "debug")

	entity.Fnc = func(b []byte, logger *logrus.Logger) error {
		return nil
	}
	assert.Nil(t, SetExecutorDefault(entity, "config.yaml"))

	entity.Fnc = func(b []byte, logger *logrus.Logger) error {
		return errors.New("ddd")
	}

	assert.Nil(t, SetExecutorDefault(entity, "config.yaml"))
}

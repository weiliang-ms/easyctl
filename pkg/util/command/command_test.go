package command

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetExecutorDefault(t *testing.T) {
	// test -c is nil
	assert.Equal(t, RunErr{}, SetExecutorDefault(Item{}))

	// test debug flag
	entity := Item{}
	entity.Cmd = &cobra.Command{}
	parentCmd := &cobra.Command{}
	parentCmd.AddCommand(entity.Cmd)
	parentParentCmd := &cobra.Command{}
	parentParentCmd.AddCommand(parentCmd)
	entity.ConfigFilePath = "config.yaml"
	assert.EqualError(t, SetExecutorDefault(entity), "flag accessed but not defined: debug")

	// test err path configfile & full in DefaultConfig
	var debug bool

	entity = Item{}
	entity.Cmd = &cobra.Command{}
	entity.DefaultConfig = []byte("ddd")
	parentCmd = &cobra.Command{}
	parentCmd.AddCommand(entity.Cmd)
	parentParentCmd = &cobra.Command{}
	parentParentCmd.AddCommand(parentCmd)
	parentParentCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug")
	entity.ConfigFilePath = "1.txt"

	assert.Errorf(t, SetExecutorDefault(entity), "open 1.txt: The system cannot find the file specified.")

}

// test panic return
func TestSetExecutorDefaultReturnErr(t *testing.T) {
	var debug bool
	entity := Item{}
	entity.Cmd = &cobra.Command{}
	parentCmd := &cobra.Command{}
	parentCmd.AddCommand(entity.Cmd)
	parentParentCmd := &cobra.Command{}
	parentParentCmd.AddCommand(parentCmd)
	parentParentCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug")

	entity.Fnc = func(item OperationItem) RunErr {
		return RunErr{}
	}
	entity.ConfigFilePath = "config.yaml"

	assert.Equal(t, RunErr{}, SetExecutorDefault(entity))

	entity.Fnc = func(item OperationItem) RunErr {
		return RunErr{Err: errors.New("ddd")}
	}

	assert.NotEqual(t, RunErr{}, SetExecutorDefault(entity))
}

// test logrus debug
func TestSetLogrusDebug(t *testing.T) {
	var debug bool
	entity := Item{}
	entity.Cmd = &cobra.Command{}
	parentCmd := &cobra.Command{}
	parentCmd.AddCommand(entity.Cmd)
	parentParentCmd := &cobra.Command{}
	parentParentCmd.AddCommand(parentCmd)
	parentParentCmd.PersistentFlags().BoolVar(&debug, "debug", true, "debug")

	entity.Fnc = func(item OperationItem) RunErr {
		return RunErr{}
	}
	entity.ConfigFilePath = "config.yaml"
	assert.Equal(t, RunErr{}, SetExecutorDefault(entity))
}

func TestRunErr_Error(t *testing.T) {
	err := RunErr{Err: errors.New("ddd"), Msg: "aaa"}
	assert.Equal(t, "ddd", err.Error())
}

package runner

import (
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"os"
	"runtime"
	"testing"
)

func TestLocalRun(t *testing.T) {
	// test error cmd
	result := LocalRun("ddd", nil)
	assert.NotNil(t, result.Err)

	// test success cmd
	if runtime.GOOS == "linux" {
		result = LocalRun("ls", nil)
		assert.Nil(t, result.Err)
	}
}

func TestParallelRun(t *testing.T) {

	server := ServerInternal{
		Host:     "10.10.10.10",
		Port:     "22",
		Password: "123",
		Username: "root",
	}

	script := "ls"

	// test shell
	executor := ExecutorInternal{
		Servers: []ServerInternal{server},
		Script:  script,
	}

	re := executor.ParallelRun()
	for v := range re {
		assert.NotNil(t, v.Err)
		assert.Equal(t, script, v.Cmd)
	}

	// test shell-script file
	script = "1.sh"
	f, _ := os.Create(script)
	f.Write([]byte("pwd"))
	f.Close()
	executor = ExecutorInternal{
		Servers: []ServerInternal{server},
		Script:  script,
	}

	re = executor.ParallelRun()
	for v := range re {
		assert.NotNil(t, v.Err)
		assert.Equal(t, "pwd", v.Cmd)
	}
	os.Remove(script)
}

func TestCompleteDefault(t *testing.T) {
	server := ServerInternal{
		Password: "123",
	}
	s := server.completeDefault()
	assert.Equal(t, "22", s.Port)
	assert.Equal(t, constant.Root, s.Username)
	assert.Equal(t, constant.RsaPrvPath, s.PrivateKeyPath)
}

func TestPublicKeyAuthFunc(t *testing.T) {
	_, err := publicKeyAuthFunc("/ss/ss")
	assert.Errorf(t, err, "ssh key file read failed open /ss/ss: The system cannot find the path specified.")

	// test not read
	_, err = publicKeyAuthFunc("~/.ssh/1.pub")
	assert.NotNil(t, err)

}

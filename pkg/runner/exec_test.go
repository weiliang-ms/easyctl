/*
	MIT License

Copyright (c) 2020 xzx.weiliang

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/
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
		UserName: "root",
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
	assert.Equal(t, constant.Root, s.UserName)
	assert.Equal(t, constant.RsaPrvPath, s.PrivateKeyPath)
}

func TestPublicKeyAuthFunc(t *testing.T) {
	_, err := publicKeyAuthFunc("/ss/ss")
	assert.Errorf(t, err, "ssh key file read failed open /ss/ss: The system cannot find the path specified.")

	// test not read
	_, err = publicKeyAuthFunc("~/.ssh/1.pub")
	assert.NotNil(t, err)

}

//func TestServerInternal_ReturnRunResult(t *testing.T) {
//	server := ServerInternal{
//		Host:     "192.168.109.160",
//		Port:     "22",
//		Password: "1",
//		UserName: "root",
//	}
//
//	re := server.ReturnRunResult(RunItem{
//		Logger: logrus.New(),
//		Cmd:    "date",
//	})
//
//	fmt.Printf("#%v", re)
//}

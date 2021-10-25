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
	"os"
	"runtime"
	"testing"
)

func TestScpErrorPathFile(t *testing.T) {
	item := ScpItem{SrcPath: "1.tt"}
	server := ServerInternal{}
	err := server.Scp(item)
	switch runtime.GOOS {
	case "windows":
		assert.EqualError(t, err, "CreateFile 1.tt: The system cannot find the file specified.")
	case "linux":
		assert.EqualError(t, err, "stat 1.tt: no such file or directory")
	}
}

func TestScpNilFile(t *testing.T) {
	f, _ := os.Create("1.txt")
	item := ScpItem{SrcPath: "1.txt"}
	server := ServerInternal{}
	err := server.Scp(item)
	assert.EqualError(t, err, "源文件: 1.txt 大小为0")
	_ = f.Close()
	_ = os.Remove("1.txt")
}

func TestConnectErr(t *testing.T) {
	f, _ := os.Create("1.txt")
	_, _ = f.Write([]byte("123"))
	item := ScpItem{SrcPath: "1.txt"}
	server := ServerInternal{}
	err := server.Scp(item)
	assert.NotNil(t, err)
	_ = f.Close()
	_ = os.Remove("1.txt")
}

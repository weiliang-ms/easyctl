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
package scp

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"io/fs"
	"testing"
)

func TestScpItem_NewSftpClient(t *testing.T) {
	item := &ScpItem{
		Servers:       nil,
		Logger:        nil,
		SftpInterface: nil,
		SftpExecutor:  SftpExecutor{},
		mock:          false,
	}

	defer func() {
		r := recover()
		if r != nil {
			require.Equal(t,
				"runtime error: invalid memory address or nil pointer dereference",
				fmt.Sprintf("%s", r))
		}
	}()

	item.NewSftpClient(runner.ServerInternal{})

}

func TestScpItem_SftpCreate(t *testing.T) {
	item := &ScpItem{
		Servers:       nil,
		Logger:        nil,
		SftpInterface: nil,
		SftpExecutor:  SftpExecutor{},
		mock:          false,
	}

	defer func() {
		r := recover()
		if r != nil {
			require.Equal(t,
				"runtime error: invalid memory address or nil pointer dereference",
				fmt.Sprintf("%s", r))
		}
	}()

	item.SftpCreate("/root/AAA.txt")

}

func TestScpItem_SftpChmod(t *testing.T) {
	item := &ScpItem{
		Servers:       nil,
		Logger:        nil,
		SftpInterface: nil,
		SftpExecutor:  SftpExecutor{},
		mock:          false,
	}

	defer func() {
		r := recover()
		if r != nil {
			require.Equal(t,
				"runtime error: invalid memory address or nil pointer dereference",
				fmt.Sprintf("%s", r))
		}
	}()

	item.SftpChmod("/root/BBB.txt", 0644)

}

func TestIOCopy64_FileNotFound(t *testing.T) {

	path := "CCC.txt"
	host := "192.168.1.1"
	logger := logrus.New()

	item := &ScpItem{
		Servers:       nil,
		Logger:        nil,
		SftpInterface: nil,
		SftpExecutor:  SftpExecutor{},
		mock:          false,
	}

	item.SrcPath = path
	item.DstPath = path

	err := item.IOCopy64(63, path, path, host, logger)

	_, ok := err.(*fs.PathError)
	require.Equal(t, true, ok)

}

// todo: mock ProxyReader

//func TestIOCopy64_FileNotFound(t *testing.T) {
//	path := "CCC.txt"
//	f , err := os.Create(path)
//	if err != nil {
//		panic(err)
//	}
//
//	host := "192.168.1.1"
//	logger := logrus.New()
//
//	item := &ScpItem{
//		Servers:       nil,
//		Logger:        nil,
//		SftpInterface: nil,
//		SftpExecutor:  SftpExecutor{},
//		mock:          false,
//	}
//
//	item.SrcPath = path
//	item.DstPath = path
//
//	err = item.IOCopy64(63, path, path, host, logger)
//	if err != nil {
//		panic(err)
//	}
//
//	f.Close()
//	err = os.Remove(path)
//	if err != nil {
//		panic(err)
//	}
//
//}

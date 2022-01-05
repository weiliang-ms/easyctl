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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/runner/scp/mocks"
	"io/fs"
	"os"
	"runtime"
	"testing"
)

func TestScpErrorPathFile(t *testing.T) {
	server := runner.ServerInternal{}
	item := &ScpItem{}
	item.SrcPath = "1.txt"
	item.ServerInternal = server

	err := Scp(item)
	switch runtime.GOOS {
	//case "windows":
	//	assert.EqualError(t, err, "CreateFile 1.tt: The system cannot find the file specified.")
	case "linux":
		assert.EqualError(t, err, "stat 1.tt: no such file or directory")
	}
}

func TestScpNilFile(t *testing.T) {
	f, _ := os.Create("222.txt")
	server := runner.ServerInternal{}
	item := &ScpItem{}
	item.SrcPath = "222.txt"
	item.ServerInternal = server
	err := Scp(item)
	assert.EqualError(t, err, "源文件: 222.txt 大小为0")
	f.Close()
	if err := os.Remove("222.txt"); err != nil {
		panic(err)
	}
}

func TestConnectErr(t *testing.T) {
	f, _ := os.Create("1.txt")
	_, _ = f.Write([]byte("123"))
	server := runner.ServerInternal{}
	item := ScpItem{}
	item.SrcPath = "1.txt"
	item.ServerInternal = server
	err := Scp(&item)
	assert.NotNil(t, err)
	f.Close()
	if err = os.Remove("1.txt"); err != nil {
		panic(err)
	}
}

//// mock
//// https://geektutu.com/post/quick-gomock.html
func TestScp(t *testing.T) {

	f, _ := os.Create("1.txt")
	_, _ = f.Write([]byte("123"))

	item := ScpItem{
		Servers: nil,
		Logger:  nil,
	}

	item.SrcPath = "1.txt"
	item.DstPath = "2.txt"
	item.Mode = 0644
	item.ServerInternal = runner.ServerInternal{}
	item.mock = true
	err := Scp(&item)

	ff, err := os.Create("2.txt")
	item.ProxyReader = ff

	if err != nil {
		fmt.Println(err)
	}

	f.Close()
	ff.Close()
	os.Remove("1.txt")
	os.Remove("2.txt")
}

func TestNewSftpProxyReader64(t *testing.T) {
	f, _ := os.Create("1.txt")
	f.WriteString("dddd")
	newSftpProxyReader64(100, "pbar", f)
	f.Close()
	os.Remove("1.txt")
}

// mock SftpInterface
func TestSftpWithProcessBar_Success_Mock(t *testing.T) {

	host := "192.168.1.1"
	size := int64(100)
	mode := fs.FileMode(0644)
	dstPath := "/root/1.txt"
	srcPath := "/root/1.txt"
	logger := logrus.New()

	mockThirdParty := &mocks.SftpInterface{}
	mockThirdParty.On("NewSftpClient", runner.ServerInternal{Host: host}).Return(nil).Once()
	mockThirdParty.On("SftpChmod", dstPath, mode).Return(nil).Once()
	mockThirdParty.On("SftpCreate", dstPath).Return(nil).Once()
	mockThirdParty.On("IOCopy64", size, srcPath, dstPath, host, logger).Return(nil).Once()

	mockScpItem := &ScpItem{
		SftpInterface: mockThirdParty,
		Logger:        logger,
		mock:          true,
	}

	mockScpItem.Mode = mode
	mockScpItem.fileSize = size
	mockScpItem.Host = host
	mockScpItem.DstPath = dstPath
	mockScpItem.SrcPath = srcPath

	err := sftpWithProcessBar(mockScpItem)
	require.Equal(t, nil, err)
}

func TestSftpWithProcessBar_NewSftpClient_Fail_Mock(t *testing.T) {

	host := "192.168.1.1"
	size := int64(100)
	mode := fs.FileMode(0644)
	srcPath := "/root/1.txt"
	dstPath := "/root/1.txt"
	logger := logrus.New()

	newSftpClientErr := fmt.Errorf("NewSftpClient Error")

	mockThirdParty := &mocks.SftpInterface{}
	mockThirdParty.On("NewSftpClient", runner.ServerInternal{Host: host}).Return(newSftpClientErr).Once()
	mockThirdParty.On("SftpChmod", dstPath, mode).Return(nil).Once()
	mockThirdParty.On("SftpCreate", dstPath).Return(nil).Once()
	mockThirdParty.On("IOCopy64", size, srcPath, dstPath, host, logger).Return(nil).Once()

	mockScpItem := &ScpItem{
		SftpInterface: mockThirdParty,
		Logger:        logger,
		mock:          true,
	}

	mockScpItem.Mode = mode
	mockScpItem.fileSize = size
	mockScpItem.Host = host
	mockScpItem.SrcPath = srcPath
	mockScpItem.DstPath = dstPath

	err := sftpWithProcessBar(mockScpItem)
	require.Equal(t, NewSftpClientErr{
		Host: host,
		Err:  newSftpClientErr,
	}, err)
}

func TestSftpWithProcessBar_SftpChmod_Fail_Mock(t *testing.T) {

	host := "192.168.1.1"
	size := int64(100)
	mode := fs.FileMode(0644)
	srcPath := "/root/1.txt"
	dstPath := "/root/1.txt"
	logger := logrus.New()

	chmodErr := fmt.Errorf("SftpChmod Error")

	mockThirdParty := &mocks.SftpInterface{}
	mockThirdParty.On("NewSftpClient", runner.ServerInternal{Host: host}).Return(nil).Once()
	mockThirdParty.On("SftpChmod", dstPath, mode).Return(chmodErr).Once()
	mockThirdParty.On("SftpCreate", dstPath).Return(nil).Once()
	mockThirdParty.On("IOCopy64", size, srcPath, dstPath, host, logger).Return(nil).Once()

	mockScpItem := &ScpItem{
		SftpInterface: mockThirdParty,
		Logger:        logger,
		mock:          true,
	}

	mockScpItem.Mode = mode
	mockScpItem.fileSize = size
	mockScpItem.Host = host
	mockScpItem.SrcPath = srcPath
	mockScpItem.DstPath = dstPath

	err := sftpWithProcessBar(mockScpItem)
	require.Equal(t, SftpChmodErr{
		Host:    host,
		DstPath: dstPath,
		Err:     chmodErr,
	}, err)
}

func TestSftpWithProcessBar_SftpCreateDst_Fail_Mock(t *testing.T) {

	host := "192.168.1.1"
	size := int64(100)
	mode := fs.FileMode(0644)
	srcPath := "/root/1.txt"
	dstPath := "/root/1.txt"
	logger := logrus.New()

	createErr := fmt.Errorf("SftpCreate Error")

	mockThirdParty := &mocks.SftpInterface{}
	mockThirdParty.On("NewSftpClient", runner.ServerInternal{Host: host}).Return(nil).Once()
	mockThirdParty.On("SftpChmod", dstPath, mode).Return(nil).Once()
	mockThirdParty.On("SftpCreate", dstPath).Return(createErr).Once()
	mockThirdParty.On("IOCopy64", size, srcPath, dstPath, host, logger).Return(nil).Once()

	mockScpItem := &ScpItem{
		SftpInterface: mockThirdParty,
		Logger:        logger,
		mock:          true,
	}

	mockScpItem.Mode = mode
	mockScpItem.fileSize = size
	mockScpItem.Host = host
	mockScpItem.SrcPath = srcPath
	mockScpItem.DstPath = dstPath

	err := sftpWithProcessBar(mockScpItem)
	require.Equal(t, SftpCreateErr{
		Host:    host,
		DstPath: dstPath,
		Err:     createErr,
	}, err)
}

func TestSftpWithProcessBar_IOCopy64_Fail_Mock(t *testing.T) {

	host := "192.168.1.1"
	size := int64(100)
	mode := fs.FileMode(0644)
	srcPath := "/root/1.txt"
	dstPath := "/root/1.txt"
	logger := logrus.New()

	ioCopy64Err := fmt.Errorf("IOCopy64 Error")

	mockThirdParty := &mocks.SftpInterface{}
	mockThirdParty.On("NewSftpClient", runner.ServerInternal{Host: host}).Return(nil).Once()
	mockThirdParty.On("SftpChmod", dstPath, mode).Return(nil).Once()
	mockThirdParty.On("SftpCreate", dstPath).Return(nil).Once()
	mockThirdParty.On("IOCopy64", size, srcPath, dstPath, host, logger).Return(ioCopy64Err).Once()

	mockScpItem := &ScpItem{
		SftpInterface: mockThirdParty,
		Logger:        logger,
		mock:          true,
	}

	mockScpItem.Mode = mode
	mockScpItem.fileSize = size
	mockScpItem.Host = host
	mockScpItem.DstPath = dstPath
	mockScpItem.SrcPath = srcPath

	err := sftpWithProcessBar(mockScpItem)
	require.Equal(t, IOCopy64Err{
		Host:    host,
		SrcPath: srcPath,
		DstPath: dstPath,
		Err:     ioCopy64Err,
	}, err)
}

func TestNewSftpClientErr_Error(t *testing.T) {
	host := "192.168.2.2"
	err := fmt.Errorf("chmod error")
	require.Equal(t, fmt.Sprintf("连接远程主机：%s失败 -> %s",
		host, err), NewSftpClientErr{
		Host: host,
		Err:  err,
	}.Error())
}

func TestSftpChmodErr_Error(t *testing.T) {
	host := "192.168.2.2"
	dstPath := "/root/abc.sh"
	err := fmt.Errorf("chmod error")
	require.Equal(t, fmt.Sprintf("修改%s:%s文件权限失败 -> %s",
		host, dstPath, err), SftpChmodErr{
		Host:    host,
		DstPath: dstPath,
		Err:     err,
	}.Error())
}

func TestSftpSftpCreateErr_Error(t *testing.T) {
	host := "192.168.2.2"
	dstPath := "/root/abc.sh"
	err := fmt.Errorf("sftp create error")
	require.Equal(t, fmt.Sprintf("创建远程主机：%s文件: %s失败 -> %s",
		host, dstPath, err), SftpCreateErr{
		Host:    host,
		DstPath: dstPath,
		Err:     err,
	}.Error())
}

func TestIOCopy64Err_Error(t *testing.T) {
	host := "192.168.2.2"
	srcPath := "/root/abc.sh"
	dstPath := "/root/abc.sh"
	err := fmt.Errorf("IOCopy64Err")
	require.Equal(t, fmt.Sprintf("拷贝文件%s至%s:%s失败 -> %s",
		srcPath, host, dstPath, err), IOCopy64Err{
		Host:    host,
		SrcPath: srcPath,
		DstPath: dstPath,
		Err:     err,
	}.Error())
}

// case: enable to connect remote server (local Test)
func TestScpLimitRate(t *testing.T) {
	s := runner.ServerInternal{
		Host:     "192.168.109.160",
		Port:     "22",
		Username: "root",
		Password: "1",
	}
	scpItem := ScpItem{
		Servers: nil,
		Logger:  logrus.New(),
	}

	path := "ddd.test"

	f, err := os.Create(path)

	if err != nil {
		panic(err)
	}

	err = f.Truncate(1024 * 1024 * 50)
	if err != nil {
		panic(err)
	}

	scpItem.SrcPath = path
	scpItem.DstPath = path

	scpItem.Mode = 0644
	scpItem.SftpExecutor.ServerInternal = s

	err = Scp(&scpItem)
	if err != nil {
		panic(err)
	}

	f.Close()
	os.Remove(path)
}

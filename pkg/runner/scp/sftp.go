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
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"io"
	"os"
	"time"
)

type SftpItem struct{}

//go:generate mockery --name=SftpInterface
type SftpInterface interface {
	NewSftpClient(server runner.ServerInternal) (err error)
	SftpCreate(path string) (err error)
	SftpChmod(dstPath string, mode os.FileMode) error
	IOCopy64(size int64, srcPath string, dstPath string, hostSign string, logger *logrus.Logger) (err error)
}

// NewSftpClient sftp客户端
func (sftp *ScpItem) NewSftpClient(server runner.ServerInternal) (err error) {
	sftp.SftpClient, err = runner.SftpConnect(server.Username, server.Password, server.Host, server.Port)
	return
}

func (sftp *ScpItem) SftpCreate(path string) (err error) {
	sftp.dstFile, err = sftp.SftpClient.Create(path)
	return
}

func (sftp *ScpItem) SftpChmod(dstPath string, mode os.FileMode) error {
	return sftp.SftpClient.Chmod(dstPath, mode)
}

func (sftp *ScpItem) IOCopy64(size int64, srcPath string, dstPath string, hostSign string, logger *logrus.Logger) (err error) {

	start := time.Now()
	sftp.srcFile, err = os.Open(srcPath)
	if err != nil {
		return err
	}
	sftp.P, sftp.ProxyReader = newSftpProxyReader64(size, hostSign, io.LimitReader(sftp.srcFile, size))

	if _, ioErr := io.Copy(sftp.dstFile, sftp.ProxyReader); ioErr != nil {
		return fmt.Errorf("传输%s:%s失败 ->%s",
			sftp.Host, sftp.DstPath, ioErr)
	}

	sftp.P.Wait()
	logger.Infof("<- %s:%s %s传输完毕...", hostSign, dstPath, time.Since(start).String())

	defer sftp.srcFile.Close()
	defer sftp.SftpClient.Close()
	defer sftp.dstFile.Close()
	defer sftp.ProxyReader.Close()
	return nil
}

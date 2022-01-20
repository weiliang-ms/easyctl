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
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/log"
	"io"
	"os"
	"sync"
	"time"
)

// ScpItem 定义跨主机拷贝属性
type ScpItem struct {
	Logger *logrus.Logger
	SftpInterface
	SftpExecutor
	SftpConnectTimeout time.Duration
	mock               bool
}

type SftpExecutor struct {
	SrcPath string
	DstPath string
	Mode    os.FileMode
	Server  runner.ServerInternal
	sync.Mutex
	SftpClient  *sftp.Client
	srcFile     *os.File
	dstFile     *sftp.File
	fileSize    int64
	ProxyReader io.ReadCloser
	P           *mpb.Progress
}

// Scp 远程写文件
func Scp(sftpItem *ScpItem) error {

	if err := fileInfo(sftpItem); err != nil {
		return err
	}

	return sftpWithProcessBar(sftpItem)
}

// ParallelScp 并发拷贝
// todo: Modify
func ParallelScp(sftpItem ScpItem, servers []runner.ServerInternal) chan error {

	ch := make(chan error, len(servers))
	wg := sync.WaitGroup{}
	wg.Add(len(servers))

	for _, s := range servers {
		go func(server runner.ServerInternal) {
			sftpItem.Lock()
			sftpItem.Server = server
			sftpItem.Unlock()
			ch <- Scp(&sftpItem)
			defer wg.Done()
		}(s)
	}
	wg.Wait()
	close(ch)
	return ch
}

type NewSftpClientErr struct {
	Host string
	Err  error
}

type SftpChmodErr struct {
	Host    string
	DstPath string
	Err     error
}

type SftpCreateErr struct {
	Host    string
	DstPath string
	Err     error
}

type IOCopy64Err struct {
	Host    string
	SrcPath string
	DstPath string
	Err     error
}

func (newSftpClientErr NewSftpClientErr) Error() string {
	return fmt.Sprintf("连接远程主机：%s失败 -> %s",
		newSftpClientErr.Host, newSftpClientErr.Err)
}

func (chmodErr SftpChmodErr) Error() string {
	return fmt.Sprintf("修改%s:%s文件权限失败 -> %s",
		chmodErr.Host, chmodErr.DstPath, chmodErr.Err)
}

func (sftpCreateErr SftpCreateErr) Error() string {
	return fmt.Sprintf("创建远程主机：%s文件: %s失败 -> %s",
		sftpCreateErr.Host, sftpCreateErr.DstPath, sftpCreateErr.Err)
}

func (ioCopy64Err IOCopy64Err) Error() string {
	return fmt.Sprintf("拷贝文件%s至%s:%s失败 -> %s",
		ioCopy64Err.SrcPath, ioCopy64Err.Host, ioCopy64Err.DstPath, ioCopy64Err.Err)
}

// todo: 优化异常处理、sftp逻辑
func sftpWithProcessBar(scpItem *ScpItem) (err error) {

	scpItem.Logger = log.SetDefault(scpItem.Logger)
	// todo: reflect
	if scpItem.SftpInterface == nil {
		scpItem.SftpInterface = new(ScpItem)
	}

	if err := scpItem.SftpInterface.NewSftpClient(scpItem.Server, scpItem.SftpConnectTimeout); err != nil {
		return NewSftpClientErr{
			Host: scpItem.Server.Host,
			Err:  err,
		}
	}

	if err := scpItem.SftpInterface.SftpCreate(scpItem.DstPath); err != nil {
		return SftpCreateErr{
			Host:    scpItem.Server.Host,
			DstPath: scpItem.DstPath,
			Err:     err,
		}
	}

	if err := scpItem.SftpInterface.SftpChmod(scpItem.DstPath, scpItem.Mode); err != nil {
		return SftpChmodErr{
			Host:    scpItem.Server.Host,
			DstPath: scpItem.DstPath,
			Err:     err,
		}
	}

	//if err := scpItem.SftpInterface.NewProxyReader(scpItem.fileSize, scpItem.SrcPath, scpItem.Host); err != nil {
	//	return NewProxyReaderErr{
	//		Host: scpItem.Host,
	//		Err:  err,
	//	}
	//}

	if err := scpItem.SftpInterface.IOCopy64(scpItem.fileSize,
		scpItem.SrcPath,
		scpItem.DstPath,
		scpItem.Server.Host,
		scpItem.Logger); err != nil {
		return IOCopy64Err{
			Host:    scpItem.Server.Host,
			SrcPath: scpItem.SrcPath,
			DstPath: scpItem.DstPath,
			Err:     err,
		}
	}

	return nil
}

func newSftpProxyReader64(size int64, name string, reader io.Reader) (*mpb.Progress, io.ReadCloser) {
	// 初始化进度条
	p := mpb.New(
		mpb.WithWidth(80),                // 进度条长度
		mpb.WithRefreshRate(time.Second), // 刷新速度
	)

	bar := p.New(size,
		// BarFillerBuilder with custom style
		mpb.BarStyle().Lbound("[").Filler("=").Tip(">").Padding("-").Rbound("|"),
		mpb.PrependDecorators(
			// display our name with one space on the right
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
			decor.CountersKibiByte("[ % .2f / % .2f ]"),
			// replace ETA decorator with "done" message, OnComplete event
			decor.NewAverageSpeed(decor.UnitKB, " % .1f", time.Now(), decor.WC{W: 6}),
			decor.OnComplete(
				decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 6}), fmt.Sprintf("done"),
			),
		),
		mpb.AppendDecorators(decor.NewPercentage("%.2f", decor.WC{W: 6})),
	)

	// create proxy reader
	return p, bar.ProxyReader(reader)
}

// todo: 检测文件合法性
func fileInfo(sftpItem *ScpItem) error {
	f, err := os.Open(sftpItem.SrcPath)
	defer f.Close()

	if err != nil {
		return err
	}

	logger := log.SetDefault(sftpItem.Logger)

	// 获取文件大小信息
	ff, _ := os.Stat(sftpItem.SrcPath)

	if ff.Size() == 0 {
		return fmt.Errorf("源文件: %s 大小为0", sftpItem.SrcPath)
	}

	sftpItem.fileSize = ff.Size()
	logger.Infof("文件大小为：%d", sftpItem.fileSize)

	return nil
}

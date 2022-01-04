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
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"io"
	"os"
	"sync"
	"time"
)

// ScpItem 定义跨主机拷贝属性
type ScpItem struct {
	Servers        []ServerInternal
	SrcPath        string
	DstPath        string
	Mode           os.FileMode
	Logger         *logrus.Logger
	ShowProcessBar bool
}

// Scp 远程写文件
func (server ServerInternal) Scp(item ScpItem) error {
	statF, err := os.Stat(item.SrcPath)

	if err != nil {
		return err
	}

	if statF.Size() == 0 {
		return fmt.Errorf("源文件: %s 大小为0", statF.Name())
	}

	// init sftp
	sftp, connetErr := sftpConnect(server.Username, server.Password, server.Host, server.Port)
	if connetErr != nil {
		return fmt.Errorf("连接远程主机：%s失败 ->%s",
			server.Host, connetErr.Error())
	}

	defer sftp.Close()

	item.Logger.Infof("-> transfer %s to %s:%s", item.SrcPath, server.Host, item.DstPath)
	dstFile, createErr := sftp.Create(item.DstPath)
	defer dstFile.Close()

	if createErr != nil {
		return fmt.Errorf("创建远程主机：%s文件: %s失败 ->%s",
			server.Host, item.DstPath, createErr.Error())
	}

	modErr := sftp.Chmod(item.DstPath, item.Mode)
	if modErr != nil {
		return fmt.Errorf("修改%s:%s文件权限失败 ->%s",
			server.Host, item.DstPath, createErr.Error())
	}

	f, _ := os.Open(item.SrcPath)
	defer f.Close()

	// 获取文件大小信息
	ff, _ := os.Stat(item.SrcPath)
	item.Logger.Infof("文件大小为：%d", ff.Size())
	reader := io.LimitReader(f, ff.Size())

	// 初始化进度条
	p := mpb.New(
		mpb.WithWidth(80),                // 进度条长度
		mpb.WithRefreshRate(time.Second), // 刷新速度
	)

	start := time.Now()

	bar := p.New(ff.Size(),
		// BarFillerBuilder with custom style
		mpb.BarStyle().Lbound("[").Filler("=").Tip(">").Padding("-").Rbound("|"),
		mpb.PrependDecorators(
			// display our name with one space on the right
			decor.Name(server.Host, decor.WC{W: len(server.Host) + 1, C: decor.DidentRight}),
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
	proxyReader := bar.ProxyReader(reader)
	_, ioErr := io.Copy(dstFile, proxyReader)

	if ioErr != nil {
		return fmt.Errorf("传输%s:%s失败 ->%s",
			server.Host, item.DstPath, createErr.Error())
	}

	p.Wait()

	defer proxyReader.Close()

	item.Logger.Infof("<- %s:%s %s传输完毕...", server.Host, item.DstPath, time.Since(start).String())
	return nil
}

// ParallelScp 并发拷贝
func ParallelScp(item ScpItem) chan error {

	ch := make(chan error, len(item.Servers))
	wg := sync.WaitGroup{}
	wg.Add(len(item.Servers))

	item.ShowProcessBar = item.Logger.Level == logrus.DebugLevel || item.ShowProcessBar

	for _, s := range item.Servers {
		go func(server ServerInternal) {
			ch <- server.Scp(item)
			defer wg.Done()
		}(s)
	}
	wg.Wait()
	close(ch)
	return ch
}

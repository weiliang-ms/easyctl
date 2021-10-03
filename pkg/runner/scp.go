package runner

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
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

	if _, err := os.Stat(item.SrcPath); err != nil {
		return err
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
		return fmt.Errorf("创建远程主机：%s文件句柄: %s失败 ->%s",
			server.Host, item.DstPath, createErr.Error())
	}

	modErr := sftp.Chmod(item.DstPath, item.Mode)
	if modErr != nil {
		return fmt.Errorf("修改%s:%s文件句柄权限失败 ->%s",
			server.Host, item.DstPath, createErr.Error())
	}

	f, _ := os.Open(item.SrcPath)
	defer f.Close()

	if item.ShowProcessBar {
		// 获取文件大小信息
		ff, _ := os.Stat(item.SrcPath)
		reader := io.LimitReader(f, ff.Size())

		// 初始化进度条
		p := mpb.New(
			mpb.WithWidth(60),                  // 进度条长度
			mpb.WithRefreshRate(1*time.Second), // 刷新速度
		)

		//
		bar := p.Add(ff.Size(),
			mpb.NewBarFiller("[=>-|"),
			mpb.PrependDecorators(
				decor.Name(fmt.Sprintf("-> %s | ", server.Host)),
				decor.CountersKibiByte("% .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.EwmaETA(decor.ET_STYLE_GO, 90),
				decor.Name(" ] "),
				decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
			),
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
	} else {
		// create proxy reader
		_, ioErr := io.Copy(dstFile, f)
		if ioErr != nil {
			return fmt.Errorf("传输%s:%s失败 ->%s",
				server.Host, item.DstPath, createErr.Error())
		}
	}
	item.Logger.Infof("<- %s:%s 传输完毕...", server.Host, item.DstPath)
	return nil
}

// ParallelScp 并发拷贝
func ParallelScp(item ScpItem) chan error {
	wg := sync.WaitGroup{}
	ch := make(chan error, len(item.Servers))
	wg.Add(len(item.Servers))

	for _, s := range item.Servers {
		go func(server ServerInternal) {
			item.ShowProcessBar = item.Logger.Level == logrus.DebugLevel || item.ShowProcessBar
			ch <- server.Scp(item)
			defer wg.Done()
		}(s)
	}

	wg.Wait()
	//close(ch)

	return ch
}

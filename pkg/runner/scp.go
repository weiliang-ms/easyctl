package runner

import (
	"fmt"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
	"io"
	"k8s.io/klog"
	"log"
	"os"
	"sync"
	"time"
)

// Scp 远程写文件
func (server ServerInternal) Scp(srcPath string, dstPath string, mode os.FileMode, showProcessBar bool) error {
	// init sftp
	sftp, connetErr := sftpConnect(server.Username, server.Password, server.Host, server.Port)
	if connetErr != nil {
		return fmt.Errorf("连接远程主机：%s失败 ->%s",
			server.Host, connetErr.Error())
	}

	defer sftp.Close()

	log.Printf("-> transfer %s to %s:%s\n", srcPath, server.Host, dstPath)
	dstFile, createErr := sftp.Create(dstPath)
	if createErr != nil {
		return fmt.Errorf("创建远程主机：%s文件句柄: %s失败 ->%s",
			server.Host, dstPath, createErr.Error())
	}
	defer dstFile.Close()

	modErr := sftp.Chmod(dstPath, mode)
	if modErr != nil {
		return fmt.Errorf("修改%s:%s文件句柄权限失败 ->%s",
			server.Host, dstPath, createErr.Error())
	}

	f, _ := os.Open(srcPath)
	defer f.Close()

	if showProcessBar {
		// 获取文件大小信息
		ff, _ := os.Stat(srcPath)
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
				server.Host, dstPath, createErr.Error())
		}

		p.Wait()

		defer proxyReader.Close()
	} else {
		// create proxy reader
		_, ioErr := io.Copy(dstFile, f)
		if ioErr != nil {
			return fmt.Errorf("传输%s:%s失败 ->%s",
				server.Host, dstPath, createErr.Error())
		}
	}

	klog.Infof("<- %s:%s 传输完毕...", server.Host, dstPath)
	return nil
}

// ParallelScp 并发拷贝
func ParallelScp(servers []ServerInternal, srcPath string, dstPath string, mode os.FileMode) chan error {
	wg := &sync.WaitGroup{}
	ch := make(chan error, len(servers))
	for _, s := range servers {
		go func(server ServerInternal) {
			wg.Add(1)
			err := server.Scp(srcPath, dstPath, mode, true)
			ch <- err
			defer wg.Done()
		}(s)
	}

	wg.Wait()
	return ch
}

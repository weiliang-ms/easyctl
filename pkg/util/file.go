package util

import (
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
	"io"
	"net/http"
	"os"
	"time"
)

func DownloadFile(srcPath, destPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	resp, err := http.Get(srcPath)
	if err != nil {
		defer out.Close()
		return err
	}
	defer resp.Body.Close()

	// 初始化进度条
	p := mpb.New(
		mpb.WithWidth(60), // 进度条长度
		//mpb.WithRefreshRate(180*time.Millisecond), // 刷新速度
		mpb.WithRefreshRate(1*time.Second),
	)

	//
	bar := p.Add(resp.ContentLength,
		mpb.NewBarFiller("[=>-|"),
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)

	// create proxy reader
	proxyReader := bar.ProxyReader(resp.Body)
	_, _ = io.Copy(out, proxyReader)

	p.Wait()

	defer out.Close()

	return nil
}

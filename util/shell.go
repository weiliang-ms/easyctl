package util

import (
	"easyctl/constant"
	"fmt"
	"sync"
)

type Shell struct {
	Cmd        string
	ServerList []Server
	Banner     Banner
}

type Banner struct {
	Symbols []string
	Msg     string
}

func PrintBanner(symbols []string, msg string) {
	fmt.Printf("%s %s...\n",
		PrintOrangeMulti(symbols), msg)
}

func (shell Shell) Shell() {
	fmt.Println(shell.Banner.Symbols)
	symbols := append(shell.Banner.Symbols, constant.Shell)
	if len(shell.ServerList) == 0 {
		PrintBanner(symbols, shell.Banner.Msg)
		ExecuteCmdAcceptResult(shell.Cmd)
	} else {
		wg := sync.WaitGroup{}
		wg.Add(len(shell.ServerList))
		for _, v := range shell.ServerList {
			PrintBanner(append(symbols, v.Host), shell.Banner.Msg)
			go v.RemoteShellParallel(shell.Cmd, &wg)
		}
		wg.Wait()
	}
}

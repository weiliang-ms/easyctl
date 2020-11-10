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

func PrintActionBanner(symbols []string, msg string) {
	fmt.Printf("%s %s...\n",
		PrintCyanMulti(symbols), msg)
}
func PrintDirectBanner(symbols []string, msg string) {
	fmt.Printf("%s %s...\n",
		PrintOrangeMulti(symbols), msg)
}

func (shell Shell) Shell() {
	var symbols []string
	if len(shell.ServerList) == 0 {
		LocalShell(shell.Banner.Msg, shell.Cmd)
	} else {
		RemoteParallelShell(symbols, shell.Banner.Msg, shell.Cmd, shell.ServerList)
	}
}

func LocalShell(msg string, cmd string) string {
	var symbols []string
	PrintDirectBanner(append(symbols, constant.LoopbackAddress), msg)
	return ExecuteCmdAcceptResult(cmd)
}

func RemoteParallelShell(symbols []string, msg string, cmd string, serverList []Server) {
	wg := sync.WaitGroup{}
	wg.Add(len(serverList))
	for _, v := range serverList {
		PrintDirectBanner(append(symbols, v.Host), msg)
		go v.RemoteShellParallel(cmd, &wg)
	}
	wg.Wait()
}

func RemoteShellAcceptResult(msg string, cmd string, server Server) string {
	var symbols []string
	PrintDirectBanner(append(symbols, server.Host), msg)
	return server.RemoteShellReturnStd(cmd)
}

func (shell Shell) ShellPrintStdout() {
	symbols := append(shell.Banner.Symbols)
	if len(shell.ServerList) == 0 {
		PrintDirectBanner(append(symbols, constant.LoopbackAddress), shell.Banner.Msg)
		ExecuteCmdPrintStd(shell.Cmd)
	} else {
		for _, v := range shell.ServerList {
			PrintDirectBanner(append(symbols, v.Host), shell.Banner.Msg)
			v.RemoteShellPrint(shell.Cmd)
		}
	}
}

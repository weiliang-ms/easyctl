package yum

//
//import (
//	"fmt"
//	"github.com/weiliang-ms/easyctl/pkg/runner"
//	"github.com/weiliang-ms/easyctl/pkg/util"
//	"log"
//	"strings"
//)
//
//// yum 安装软件实例
//type Installer struct {
//	Server   runner.Server // 主机信息
//	Software []string      // 软件集合
//}
//
//// yum 安装结果
//type InstallResult struct {
//	Host     string `table:"主机地址"`
//	Software string `table:"安装软件"`
//	Status   string `table:"安装状态"`
//}
//
//// soft 检测实例
//type Detect struct {
//	Server   []runner.Server
//	Software []string
//}
//
//// yum 检测结果
//type DetectInfo struct {
//	Host   string `table:"主机地址"`
//	Status string `table:"检测状态"`
//	Msg    string `table:"错误信息"`
//	Time   string `table:"检测时间"`
//}
//
//// 检测结果
//type DetectResult struct {
//	Info      []DetectInfo
//	Installer []Installer
//}
//
//var yumInstallCmd = "yum install -y"
//
//// 返回yum异常主机列表
//func DetectYum(server []runner.Server) (badYumServers []runner.Server) {
//	cmd := "yum repolist 2>/dev/null |grep repolist|awk '{print $2}'"
//	for _, v := range server {
//		if v.Host == util.Localhost {
//			re := runner.Shell(cmd)
//			if re.ExitCode != 0 || strings.Trim(re.StdOut, "\n") == "0" {
//				badYumServers = append(badYumServers, v)
//			}
//		} else {
//			re := v.RemoteShell(cmd)
//			if re.Code != 0 || strings.Trim(re.StdOut, "\n") == "0" {
//				badYumServers = append(badYumServers, v)
//			}
//		}
//	}
//
//	return badYumServers
//}
//
//// 返回yum异常主机列表
//func Install(ins []Installer) (errServers []InstallResult) {
//
//	for _, i := range ins {
//		if i.Server.Host == util.Localhost {
//			return localInstall(i.Software)
//		} else {
//			errServers = append(errServers, remoteInstall(i.Software, i.Server))
//		}
//	}
//
//	return
//}
//
//func localInstall(soft []string) (errServers []InstallResult) {
//	for _, s := range soft {
//		shell := fmt.Sprintf("%s %s", yumInstallCmd, s)
//		re := runner.Shell(shell)
//
//		if re.ExitCode != 0 {
//			return
//		}
//
//		errInstance := InstallResult{
//			Host:     util.Localhost,
//			Status:   util.Fail,
//			Software: s,
//		}
//
//		errServers = append(errServers, errInstance)
//
//	}
//
//	return errServers
//}
//
//func remoteInstall(soft []string, server runner.Server) (instance InstallResult) {
//	for _, s := range soft {
//		shell := fmt.Sprintf("%s %s", yumInstallCmd, s)
//		re := server.RemoteShell(shell)
//		if re.Code != 0 && re.StdOut == "0" {
//			instance = InstallResult{
//				Host:     server.Host,
//				Status:   util.Fail,
//				Software: s,
//			}
//
//		}
//
//	}
//
//	return
//}
//
//func DetectAllNodes(instance Detect) *DetectResult {
//
//	re := &DetectResult{}
//
//	if len(instance.Server) == 0 {
//		info, installer := detectLocalNode(instance.Software)
//		if info == nil || installer == nil {
//			return nil
//		}
//		re.Info = append(re.Info, *info)
//		re.Installer = append(re.Installer, *installer)
//		return re
//	}
//
//	for _, v := range instance.Server {
//
//		info, installer := detectOnNode(v, instance.Software)
//
//		if info == nil || installer == nil {
//			return nil
//		}
//
//		re.Info = append(re.Info, *info)
//		re.Installer = append(re.Installer, *installer)
//	}
//
//	return re
//}
//
//// 本地节点检测
//func detectLocalNode(softList []string) (*DetectInfo, *Installer) {
//
//	var detectFailSoft []string
//	log.Printf("检测本机依赖 %v...\n", softList)
//
//	for _, software := range softList {
//		re := runner.ShortShell(fmt.Sprintf("sudo -E /bin/sh -c \"yum -q install -y %s &>/dev/null || rpm -qa|grep %s &>/dev/null\"", software, software))
//		if re.ExitCode != 0 {
//			detectFailSoft = append(detectFailSoft, software)
//		}
//	}
//
//	if len(detectFailSoft) > 0 {
//		shellRe := runner.ShortShell("date +\"%Z %H:%M:%S\"")
//		re := &DetectInfo{}
//		needed := &Installer{}
//		re.Msg = fmt.Sprintf("%s检测失败...", util.AppendStringFromSlice(detectFailSoft, ","))
//		re.Host = util.Localhost
//		re.Status = util.Fail
//		re.Time = strings.TrimSpace(shellRe.StdOut)
//
//		needed.Server = runner.Server{
//			Host: util.Localhost,
//		}
//
//		needed.Software = detectFailSoft
//
//		return re, needed
//	}
//
//	return nil, nil
//}
//
//// 按节点检测
//func detectOnNode(server runner.Server, software []string) (*DetectInfo, *Installer) {
//
//	log.Printf("检测节点：%s\n", server.Host)
//	var detectFailSoft []string
//
//	for _, software := range software {
//		re := server.RemoteShell(fmt.Sprintf("sudo -E /bin/sh -c \"yum -q install -y %s &>/dev/null || rpm -qa|grep %s &>/dev/null\"", software, software))
//		if re.Code != 0 {
//			detectFailSoft = append(detectFailSoft, software)
//		}
//	}
//
//	if len(detectFailSoft) > 0 {
//		shellRe := server.RemoteShell("date +\"%Z %H:%M:%S\"")
//		info := &DetectInfo{}
//		instance := &Installer{}
//
//		info.Msg = fmt.Sprintf("%s检测失败...", util.AppendStringFromSlice(detectFailSoft, ","))
//		info.Host = server.Host
//		info.Status = util.Fail
//		info.Time = strings.TrimSpace(shellRe.StdOut)
//
//		instance.Server = server
//		instance.Software = detectFailSoft
//
//		return info, instance
//	}
//
//	return nil, nil
//}

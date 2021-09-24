package install

//
//import (
//	"bufio"
//	_ "embed"
//	"fmt"
//	"github.com/lithammer/dedent"
//	"github.com/pkg/errors"
//	"github.com/weiliang-ms/easyctl/pkg/runner"
//	"github.com/weiliang-ms/easyctl/pkg/util"
//	"log"
//	"os"
//	"strings"
//	"sync"
//	"text/template"
//)
//
//var dockerInstallTmpl = template.Must(template.New("dockerInstallTmpl").Parse(
//	dedent.Dedent(
//		`
//#!/bin/bash
//package_name=docker-ce
//
//cat > /etc/yum.repos.d/"$package_name".repo <<EOF
//[$package_name]
//name=[$package_name]-repo
//baseurl=file:///yum/data/"$package_name"
//gpgcheck=0
//enabled=1
//EOF
//
//# 解压缩
//mkdir -p /yum
//tar zxvf {{.PackagePath}} -C /yum/
//
//# 安装
//sudo sed -i "/net.ipv4.ip_forward/d" /etc/sysctl.conf
//sudo echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
//sudo sysctl -p
//
//cat > /etc/systemd/system/docker.service <<EOF
//[Unit]
//Description=Docker Application Container Engine
//Documentation=https://docs.docker.com
//After=network-online.target firewalld.service
//Wants=network-online.target
//[Service]
//Type=notify
//# the default is not to use systemd for cgroups because the delegate issues still
//# exists and systemd currently does not support the cgroup feature set required
//# for containers run by docker
//ExecStart=/usr/bin/dockerd -H tcp://127.0.0.1:4243 -H unix:///var/run/docker.sock
//ExecReload=/bin/kill -s HUP
//# Having non-zero Limit*s causes performance problems due to accounting overhead
//# in the kernel. We recommend using cgroups to do container-local accounting.
//LimitNOFILE=infinity
//LimitNPROC=infinity
//LimitCORE=infinity
//# Uncomment TasksMax if your systemd version supports it.
//# Only systemd 226 and above support this version.
//#TasksMax=infinity
//TimeoutStartSec=0
//# set delegate yes so that systemd does not reset the cgroups of docker containers
//Delegate=yes
//# kill only the docker process, not all processes in the cgroup
//KillMode=process
//# restart the docker process if it exits prematurely
//Restart=on-failure
//StartLimitBurst=3
//StartLimitInterval=60s
//[Install]
//WantedBy=multi-user.target
//EOF
//
//mkdir -p /etc/docker
//sudo tee /etc/docker/daemon.json <<EOF
//{
//  "log-opts": {
//    "max-size": "5m",
//    "max-file":"3"
//  },
//  "userland-proxy": false,
//  "live-restore": true,
//  "default-ulimits": {
//    "nofile": {
//      "Hard": 65535,
//      "Name": "nofile",
//      "Soft": 65535
//    }
//  },
//  "default-address-pools": [
//    {
//      "base": "172.80.0.0/16",
//      "size": 24
//    },
//    {
//      "base": "172.90.0.0/16",
//      "size": 24
//    }
//  ],
//  {{- if .DataDir }}
//  "data-root": "{{ .DataDir }}",
//  {{- end}}
//  "live-restore": true
//}
//EOF
//
//systemctl daemon-reload
//systemctl enable docker --now
//systemctl restart docker
//`)))
//
////go:embed asset/install_offline_docker.sh
//var dockerScript []byte
//
//// Docker 单机本地离线
//func Docker(i runner.Installer) {
//	i.Cmd = fmt.Sprintf(string(dockerScript))
//
//	// 截取离线安装介质名称
//	strSlice := strings.Split(i.OfflineFilePath, "/")
//	if len(strSlice) == 0 {
//		log.Fatalf("获取文件名失败...")
//	}
//	i.FileName = strSlice[len(strSlice)-1]
//
//	if i.ServerListPath != "" {
//		re := runner.ParseServerList(i.ServerListPath, runner.DockerServerList{})
//		list := re.Docker.Attribute.Servers
//		if i.Offline {
//			offlineRemote(i, list)
//		}
//	} else {
//		localOfflineInstall(i)
//	}
//}
//
//func localOfflineInstall(i runner.Installer) {
//
//	// todo: 优化关闭selinux
//	re := runner.Shell("getenforce")
//	if re.StdOut != "Disabled" {
//
//		runner.Shell("setenforce 0 && sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config")
//
//		reader := bufio.NewReader(os.Stdin)
//		input, err := confirm(reader)
//		if err != nil {
//			panic(err)
//		}
//
//		if input == "no" {
//			os.Exit(0)
//		}
//
//		runner.Shell("reboot")
//
//	}
//
//	// 安装
//	log.Println("开始本机安装docker-ce...")
//	content, _ := util.Render(dockerInstallTmpl, util.Data{
//		"PackagePath": i.OfflineFilePath,
//		"DataDir":     i.DataDir,
//	})
//
//	// 创建数据目录
//	// todo: 校验数据目录合法性
//	log.Printf("创建数据目录: %s合法性...\n", i.DataDir)
//	legal, detectErr := dirLegitimacyTest(i.DataDir)
//	if !legal {
//		log.Fatalln(errors.Wrap(detectErr, "数据目录非法，请变更!").Error())
//	}
//	log.Printf("数据目录: %s合法性检测通过...\n", i.DataDir)
//
//	err := os.MkdirAll(i.DataDir, 0755)
//	if err != nil {
//		log.Fatalf("创建数据目录: %s失败...\n", i.DataDir)
//	}
//
//	runner.Shell(content)
//	runner.Shell("docker version")
//}
//
//func offlineRemote(i runner.Installer, list []runner.Server) {
//	var wg sync.WaitGroup
//
//	ch := make(chan runner.ShellResult, len(list))
//	// 拷贝文件
//	dstPath := fmt.Sprintf("/tmp/%s", i.FileName)
//
//	// 生成本地临时文件
//	for _, v := range list {
//		runner.ScpFile(i.OfflineFilePath, dstPath, v, 0755)
//		log.Println("-> transfer done ...")
//	}
//
//	// 并行
//	log.Println("-> 批量安装...")
//	for _, v := range list {
//		wg.Add(1)
//		go func(server runner.Server) {
//			defer wg.Done()
//			re := server.RemoteShell(i.Cmd)
//			ch <- re
//		}(v)
//	}
//
//	wg.Wait()
//	close(ch)
//
//	// ch -> slice
//	var as []runner.ShellResult
//	for target := range ch {
//		as = append(as, target)
//	}
//}
//
//func confirm(reader *bufio.Reader) (string, error) {
//	for {
//		fmt.Printf("selinux已关闭，请问是否关闭并重启主机生效? [yes/no]: ")
//		input, err := reader.ReadString('\n')
//		if err != nil {
//			return "", err
//		}
//		input = strings.TrimSpace(input)
//
//		if input != "" && (input == "yes" || input == "no") {
//			return input, nil
//		}
//	}
//}
//
//// 目录合法性检测
//func dirLegitimacyTest(dirPath string) (re bool, err error) {
//
//	riskParentDirs := []string{
//		"bin",
//		"boot",
//		"dev",
//		"etc",
//		"lib",
//		"lib64",
//		"proc",
//		"root",
//		"sbin",
//		"sys",
//		"srv",
//		"tmp",
//	}
//
//	log.Printf("检测目录：%s合法性\n", dirPath)
//	splits := strings.Split(dirPath, "/")
//
//	if dirPath == "" {
//		return false, errors.New("目录路径不能为空")
//	}
//
//	if len(dirPath) >= 1 && dirPath[0:1] != "/" {
//		return false, errors.New("目录路径必须为以\"/\"起始的绝对路径")
//	}
//
//	if dirPath == "/" {
//		return false, errors.New("目录路径不能为根路径")
//	}
//
//	log.Println("检测父级目录合法性...")
//
//	if len(splits) == 2 && util.SliceContain(riskParentDirs, splits[1]) {
//		return false, errors.New(fmt.Sprintf("目录：%s为系统级目录!!!", dirPath))
//	}
//
//	return true, nil
//}

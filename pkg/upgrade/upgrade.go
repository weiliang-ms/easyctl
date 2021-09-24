package upgrade

//
//import (
//	"bufio"
//	"fmt"
//	"github.com/modood/table"
//	"github.com/weiliang-ms/easyctl/pkg/runner"
//	"github.com/weiliang-ms/easyctl/pkg/runner/yum"
//	"github.com/weiliang-ms/easyctl/pkg/util"
//	"io/ioutil"
//	"log"
//	"net/url"
//	"os"
//	"strings"
//	"sync"
//)
//
//// todo: 简化结构体
//
//// 执行器执行结果
//type ActuatorExecRe struct {
//	Host     string `table:"主机IP"`
//	Action   string `table:"执行操作"`
//	Status   string `table:"执行状态"`
//	ErrorMsg string `table:"执行结果"`
//}
//
//// 执行器执行结果
//type ShellExecRe struct {
//	Host   string `table:"主机IP"`
//	Action string `table:"执行操作"`
//	Status string `table:"执行状态"`
//	Shell  string `table:"执行语句"`
//}
//
//type ExecShellObject struct {
//	Server runner.Server
//	Shell  string
//}
//
//var Status string
//
//type Actuator struct {
//	ServerListFile   string
//	Cmd              string
//	InstallDepCmd    string // 安装依赖命令
//	FilePath         string
//	FileName         string
//	ServerList       runner.CommonServerList
//	Err              error
//	ExitCode         int
//	DependenciesList []string // 依赖列表
//	NeedShell        []ExecShellObject
//	OpensslDir       string // 本地openssl安装目录
//}
//
//func init() {
//	Status = util.Fail
//}
//
//func (ac *Actuator) download() *Actuator {
//
//	// 判断是否为 url 类型
//	_, err := url.ParseRequestURI(ac.FilePath)
//	requestUrl := ac.FilePath
//	if err == nil {
//		strSlice := strings.Split(ac.FilePath, "/")
//		if len(strSlice) > 0 {
//			ac.FilePath = strSlice[len(strSlice)-1]
//		}
//		log.Println(fmt.Sprintf("检测到文件类型为url，尝试下载: %s...", requestUrl))
//		err := util.DownloadFile(requestUrl, ac.FilePath)
//		if err != nil {
//			log.Fatalln(err.Error())
//		}
//	} else {
//		_, err := os.Stat(ac.FilePath)
//		if err != nil {
//			log.Fatalf(fmt.Sprintf("检测到：%s异常 -> %s", ac.FilePath, err.Error()))
//		}
//	}
//
//	return ac
//}
//
//func (ac *Actuator) parseServerList() *Actuator {
//
//	if ac.ServerListFile != "" {
//		serverList := runner.ParseCommonServerList(ac.ServerListFile)
//		ac.ServerList = serverList
//	}
//
//	return ac
//}
//
//// 编译/运行时 依赖检测
//func (ac *Actuator) detect() *Actuator {
//
//	instance := yum.Detect{
//		Server:   ac.ServerList.Server,
//		Software: ac.DependenciesList,
//	}
//
//	// 尝试安装依赖
//	ac.installDep(yum.DetectAllNodes(instance))
//
//	return ac
//}
//
//// 安装依赖
//func (ac *Actuator) installDep(result *yum.DetectResult) *Actuator {
//
//	if result == nil {
//		log.Println("依赖检测通过...")
//		return ac
//	}
//
//	log.Println("依赖检测未通过，检测结果如下...")
//	table.OutputA(result.Info)
//
//	log.Println("检测yum可用性...")
//	var s []runner.Server
//	for _, v := range result.Installer {
//		s = append(s, v.Server)
//	}
//	badYumServer := yum.DetectYum(s)
//	if badYumServer != nil {
//		log.Println("以下主机yum不可用，请调整...")
//		table.OutputA(badYumServer)
//		os.Exit(0)
//	}
//	log.Println("yum可用性检测通过...")
//
//	util.CountDown(5, "尝试安装依赖")
//	yum.Install(result.Installer)
//
//	return ac
//}
//
////
//func (ac *Actuator) Confirm(results interface{}) *Actuator {
//	table.OutputA(results)
//	reader := bufio.NewReader(os.Stdin)
//Loop:
//	for {
//		fmt.Printf("检测到需要安装以上依赖，是否? [yes/no]: ")
//		input, err := reader.ReadString('\n')
//		if err != nil {
//			log.Fatal(err)
//		}
//		input = strings.TrimSpace(strings.ToLower(input))
//
//		switch input {
//		case "yes":
//			break Loop
//		case "no":
//			os.Exit(0)
//		default:
//			continue
//		}
//	}
//
//	return ac
//}
//
//// todo:// 并发结果集返回与过程实时输出...
///*
//close没有make的chan会引起panic
//close以后不能再写入，写入会出现panic
//close之后可以读取，无缓冲chan读取返回0值和false，有缓冲chan可以继续读取，返回的都是chan中数据和true，直到读取完所有队列中的数据。
//重复close会引起panic
//只读chan不能close
//不close chan也是可以的，当没有被引用时系统会自动垃圾回收
//*/
//func (ac *Actuator) execute(action string, ignoreErrCode int) {
//
//	var result []ActuatorExecRe
//	log.Println(action)
//
//	ch := make(chan ActuatorExecRe, len(ac.ServerList.Server))
//
//	var wg sync.WaitGroup
//	if len(ac.ServerList.Server) <= 0 {
//		re := runner.Shell(ac.Cmd)
//		ac.ExitCode = re.ExitCode
//		result = append(result, packageRe(util.Localhost, action, re.StdErr, ignoreErrCode, re.ExitCode))
//	} else {
//		for _, v := range ac.ServerList.Server {
//			wg.Add(1)
//			go func(actionMsg string, ignoreCode int, server runner.Server) {
//				re := server.RemoteShell(ac.Cmd)
//				ioutil.WriteFile(fmt.Sprintf("%s.out", server.Host), []byte(re.StdOut), 0644)
//				ch <- packageRe(server.Host, actionMsg, re.StdErrMsg, ignoreCode, re.Code)
//				wg.Done()
//			}(action, ignoreErrCode, v)
//		}
//
//	}
//
//	defer func() {
//		close(ch)
//		for msg := range ch {
//			result = append(result, msg)
//		}
//		table.OutputA(result)
//	}()
//
//	wg.Wait()
//}
//
//// 分发文件
//func (ac *Actuator) handoutFile() *Actuator {
//
//	log.Println("分发介质文件...")
//	if len(ac.ServerList.Server) <= 0 {
//		fmt.Println("返回了...")
//		return ac
//	}
//	localPath := ac.FilePath
//	ac.FilePath = fmt.Sprintf("/tmp/%s", ac.FilePath)
//	for _, v := range ac.ServerList.Server {
//		runner.ScpFile(localPath, ac.FilePath, v, 0644)
//	}
//	return ac
//}
//
//func packageRe(host, action, errMsg string, ignoreErrCode int, code int) ActuatorExecRe {
//
//	if code == 0 || code == ignoreErrCode {
//		Status = util.Success
//	}
//
//	return ActuatorExecRe{
//		Host:     host,
//		Action:   action,
//		Status:   Status,
//		ErrorMsg: errMsg,
//	}
//}

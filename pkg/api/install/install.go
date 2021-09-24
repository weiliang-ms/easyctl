package install

//
//import (
//	"fmt"
//	"github.com/weiliang-ms/easyctl/pkg/runner"
//	"github.com/weiliang-ms/easyctl/pkg/util"
//	"gopkg.in/yaml.v2"
//	"log"
//	"net"
//	"os"
//	"strconv"
//	"strings"
//)
//
//// Dependency DependencyDetect 依赖检测实体
//type Dependency struct {
//	Needs map[string]struct {
//		DetectShell  string // 检测shell
//		InstallShell string // 安装shell
//	}
//	Server runner.Server
//}
//
//type DetectResult struct {
//	Host   string
//	Reason map[string]string
//	Status DetectStatus
//}
//
//type DetectStatus string
//
//// 执行器执行结果
//type ExecutorResult struct {
//	Host     string `table:"主机IP"`
//	Action   string `table:"执行操作"`
//	Status   string `table:"执行状态"`
//	ErrorMsg string `table:"错误信息"`
//}
//
//const (
//	YumErrMsg = "yum不可用..."
//	http      = "http"
//	https     = "https"
//	ftp       = "ftp"
//
//	DetectFail    DetectStatus = "依赖检测未通过"
//	DetectSuccess DetectStatus = "依赖检测通过"
//)
//
//// HaproxyMeta 安装配置haproxy元数据
//type HaproxyMeta struct {
//	Haproxy struct {
//		FilePath string          `yaml:"file"`
//		Server   []runner.Server `yaml:"server"`
//		Excludes []string        `yaml:"excludes,flow"`
//		Balance  []struct {
//			Name       string   `yaml:"name"`
//			ListenPort int      `yaml:"listen-port"`
//			Endpoint   []string `yaml:"endpoint"`
//		} `yaml:"balance"`
//	} `yaml:"haproxy"`
//}
//
//// ParseConfig 解析yaml配置
//func ParseConfig(content []byte, v interface{}) (interface{}, error) {
//
//	t := convert(v)
//
//	err := yaml.Unmarshal(content, t)
//	if err != nil {
//		return nil, err
//	}
//
//	return t, nil
//}
//
//// ParseServer 解析server主机列表
//func ParseServer(excludes []string, s runner.Server) []runner.Server {
//	serverList := []runner.Server{}
//
//	contain := strings.Contains(s.Host, "[") && strings.Contains(s.Host, "]") && strings.Contains(s.Host, ":")
//
//	if address := net.ParseIP(s.Host); address == nil && !contain {
//		log.Fatalln("server地址信息非法，无法解析请检查...")
//	}
//
//	if strings.Contains(s.Host, "[") {
//		log.Println("检测到配置文件中含有IP段，开始解析组装...")
//		//192.168.235.
//		baseAddress := strings.Split(s.Host, "[")[0]
//		log.Printf("解析到IP子网网段为：%s...\n", baseAddress)
//
//		// 1:3] -> 1:3
//		ipRange := strings.Split(strings.Split(s.Host, "[")[1], "]")[0]
//		log.Printf("解析到IP区间为：%s...\n", ipRange)
//
//		// 1:3 -> 1
//		begin := strings.Split(ipRange, ":")[0]
//		log.Printf("解析到起始IP为：%s...\n", fmt.Sprintf("%s%s", baseAddress, begin))
//
//		// 1:3 -> 3
//		end := strings.Split(ipRange, ":")[1]
//		log.Printf("解析到末尾IP为：%s...\n", fmt.Sprintf("%s%s", baseAddress, end))
//
//		// string -> int
//		beginIndex, _ := strconv.Atoi(begin)
//		endIndex, _ := strconv.Atoi(end)
//
//		for i := beginIndex; i <= endIndex; i++ {
//			server := runner.Server{
//				Host:           fmt.Sprintf("%s%d", baseAddress, i),
//				Port:           s.Port,
//				Username:       s.Username,
//				Password:       s.Password,
//				PrivateKeyPath: s.PrivateKeyPath,
//			}
//
//			if !util.SliceContain(excludes, server.Host) {
//				serverList = append(serverList, server)
//			}
//		}
//	} else {
//		serverList = append(serverList, s)
//	}
//
//	return serverList
//}
//
//// Download 下载安装介质
//func Download(filePath string) (string, error) {
//	// 判断是否为 url 类型
//	var preservePath string
//	hasPrefix := strings.HasPrefix(filePath, http) || strings.HasPrefix(filePath, https) || strings.HasPrefix(filePath, ftp)
//	if hasPrefix {
//		url := filePath
//		strSlice := strings.Split(url, "/")
//		if len(strSlice) > 0 {
//			preservePath = strSlice[len(strSlice)-1]
//		}
//
//		log.Println(fmt.Sprintf("检测到文件类型为url，尝试下载: %s...", filePath))
//		return preservePath, util.DownloadFile(filePath, preservePath)
//
//	} else {
//		if _, err := os.Stat(filePath); err != nil {
//			return filePath, err
//		}
//	}
//
//	return filePath, nil
//}
//
//// DetectDep 检测依赖
//func (d Dependency) DetectDep() DetectResult {
//	re := DetectResult{}
//	re.Reason = make(map[string]string)
//	for k, v := range d.Needs {
//		if r := d.Server.RemoteShell(v.DetectShell); r.Code != 0 {
//			log.Printf("%s尝试安装%s\n", d.Server.Host, k)
//			if err := d.Server.InstallSoft(v.InstallShell); err != nil {
//				re.Reason[k] = YumErrMsg
//				re.Status = DetectFail
//				re.Host = d.Server.Host
//			} else {
//				re.Host = d.Server.Host
//				re.Reason[k] = ""
//				re.Status = DetectSuccess
//			}
//		}
//		log.Printf("%s检测%s通过", d.Server.Host, k)
//	}
//
//	return re
//}
//
//// HandFile 分发安装介质
//func HandFile(serverList []runner.Server, filepath string) error {
//
//	// 截取文件名称
//	filename := subFilename(filepath)
//
//	// 生成scp拷贝文件目标路径
//	dstPath := fmt.Sprintf("/tmp/%s", filename)
//
//	// scp至各节点
//	for _, v := range serverList {
//		err := v.Scp(filepath, dstPath, 0755)
//		if err != nil {
//			return err
//		}
//		log.Printf("<- transfer %s to %s:%s done...\n", filepath, v.Host, dstPath)
//	}
//
//	return nil
//}
//
//// Execute todo: 并发执行
//func Execute(servers []runner.Server, action, cmd string) []ExecutorResult {
//
//	var result []ExecutorResult
//	if len(servers) <= 0 {
//		re := runner.Shell(cmd)
//		result = append(result, packageRe(util.Localhost, action, re.StdErr, re.ExitCode))
//	} else {
//		for _, v := range servers {
//			go v.RemoteShell(cmd)
//			//re := v.RemoteShell(cmd)
//			//result = append(result, packageRe(v.Host, action, re.StdErrMsg, re.Code))
//		}
//	}
//	return result
//}
//
//func packageRe(host, action, errMsg string, code int) ExecutorResult {
//
//	var status string
//
//	if code == 0 {
//		status = util.Success
//	} else {
//		status = util.Fail
//	}
//	return ExecutorResult{
//		Host:     host,
//		Action:   action,
//		Status:   status,
//		ErrorMsg: errMsg,
//	}
//}
//
//// 类型转换
//func convert(v interface{}) interface{} {
//	switch v.(type) {
//	case HaproxyMeta:
//		return &HaproxyMeta{}
//	default:
//		return nil
//	}
//}
//
//func subFilename(path string) string {
//
//	slice := util.SubSlash(path)
//	if len(slice) > 1 {
//		return slice[len(slice)-1]
//	}
//
//	return path
//}

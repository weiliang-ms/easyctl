package manager

//
//import (
//	"encoding/json"
//	"github.com/pkg/errors"
//	log "github.com/sirupsen/logrus"
//	"github.com/weiliang-ms/easyctl/pkg/util"
//	"github.com/weiliang-ms/easyctl/pkg/util/ssh"
//	"gopkg.in/yaml.v2"
//	"os"
//	"sync"
//	"time"
//)
//
//const (
//	DefaultCon = 10
//	Timeout    = 120
//)
//
//type Runner struct {
//	Conn   ssh.Connection
//	Debug  bool
//	Server *ssh.Server
//	Index  int
//}
//
//type ExecResult struct {
//	ExitCode int
//	StdErr   string
//	StdOut   string
//}
//
//type ShellResult struct {
//	Host    string
//	Cmd     string
//	Code    int
//	Status  string
//	Content string
//}
//
//type Manager struct {
//	Logger         log.FieldLogger
//	ServerListFile string
//	Servers        []ssh.Server
//	Connector      *ssh.Dialer
//	Runner         *Runner
//	Offline        bool
//	Cmd            string
//	FileName       string
//	Debug          bool
//	OfflineFile    string
//	MediaType      MediaType
//}
//
//type MediaType int // 软件安装类型
//
//const (
//	_ MediaType = iota // 根据iota特性定义枚举类型常量
//	YUM
//	RPM
//	BINARY
//)
//
//type Task struct {
//	Task   func(*Manager) error
//	ErrMsg string
//}
//
//type NodeTask func(mgr *Manager, server *ssh.Server) error
//
//func (mgr *Manager) Copy() *Manager {
//	newManager := *mgr
//	return &newManager
//}
//
//func (t *Task) Run(mgr *Manager) error {
//	backoff := util.Backoff{
//		Steps:    1,
//		Duration: 5 * time.Second,
//		Factor:   2.0,
//	}
//
//	var lastErr error
//	err := util.ExponentialBackoff(backoff, func() (bool, error) {
//		lastErr = t.Task(mgr)
//		if lastErr != nil {
//			mgr.Logger.Warn("Task failed ...")
//			if mgr.Debug {
//				mgr.Logger.Warnf("error: %s", lastErr)
//			}
//			return false, nil
//		}
//		return true, nil
//	})
//	if err == util.ErrWaitTimeout {
//		err = lastErr
//	}
//	return err
//}
//
//func (mgr *Manager) runTask(server *ssh.Server, task NodeTask, index int) error {
//	var (
//		err  error
//		conn ssh.Connection
//	)
//
//	conn, err = mgr.Connector.Connect(*server)
//	if err != nil {
//		return errors.Wrapf(err, "Failed to connect to %s", server.Host)
//	}
//
//	mgr.Runner = &Runner{
//		Conn:   conn,
//		Debug:  mgr.Debug,
//		Server: server,
//		Index:  index,
//	}
//
//	return task(mgr, server)
//}
//
//func (mgr *Manager) RunTaskOnAllNodes(task NodeTask, parallel bool) error {
//	if err := mgr.RunTaskOnNodes(mgr.Servers, task, parallel); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (mgr *Manager) RunTaskOnNodes(nodes []ssh.Server, task NodeTask, parallel bool) error {
//	var err error
//	hasErrors := false
//
//	wg := &sync.WaitGroup{}
//	result := make(chan string)
//	ccons := make(chan struct{}, DefaultCon)
//	defer close(result)
//	defer close(ccons)
//	hostNum := len(nodes)
//
//	if parallel {
//		go func(result chan string, ccons chan struct{}) {
//			for i := 0; i < hostNum; i++ {
//				select {
//				case <-result:
//				case <-time.After(time.Minute * Timeout):
//					mgr.Logger.Fatalf("Execute task timeout, Timeout=%ds", Timeout)
//				}
//				wg.Done()
//				<-ccons
//			}
//		}(result, ccons)
//	}
//
//	for i := range nodes {
//		mgr := mgr.Copy()
//		mgr.Logger = mgr.Logger.WithField("node", nodes[i].Host)
//
//		if parallel {
//			ccons <- struct{}{}
//			wg.Add(1)
//			go func(mgr *Manager, node *ssh.Server, result chan string, index int) {
//				err = mgr.runTask(node, task, index)
//				if err != nil {
//					mgr.Logger.Error(err)
//					hasErrors = true
//				}
//				result <- "done"
//			}(mgr, &nodes[i], result, i)
//		} else {
//			err = mgr.runTask(&nodes[i], task, i)
//			if err != nil {
//				break
//			}
//		}
//	}
//
//	wg.Wait()
//
//	if hasErrors {
//		err = errors.New("interrupted by error")
//	}
//
//	return err
//}
//
//func RemoteShell(cmd string, server ssh.Server) ShellResult {
//
//	var result ShellResult
//	if len(cmd) < 60 {
//		result.Cmd = cmd
//	} else {
//		result.Cmd = "built-in shell"
//	}
//
//	result.Host = server.Host
//	log.Printf("-> [%s] shell => %s", server.Host, cmd)
//
//	session, conErr := server.SSHConnect()
//	if conErr != nil {
//		log.Fatal("连接失败...")
//	}
//
//	defer session.Close()
//	combo, runErr := session.CombinedOutput(cmd)
//
//	if runErr != nil && runErr.Error() == "Process exited with status 1" {
//		result.Code = 1
//		return result
//	}
//
//	result.Content = string(combo)
//
//	return result
//}
//
//func (mgr *Manager) ParseServerList() error {
//
//	var decodeErr, marshalErr error
//	var servers ssh.Servers
//
//	f, err := os.Open(mgr.ServerListFile)
//	if err != nil {
//		return err
//	}
//
//	decodeErr = yaml.NewDecoder(f).Decode(&servers)
//	_, marshalErr = json.Marshal(&servers)
//
//	if decodeErr != nil {
//		return err
//	}
//
//	if marshalErr != nil {
//		return err
//	}
//
//	mgr.Servers = servers.Server
//
//	return nil
//}

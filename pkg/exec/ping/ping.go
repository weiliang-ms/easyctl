package ping

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/file"
	"github.com/weiliang-ms/easyctl/pkg/util/log"
	strings2 "github.com/weiliang-ms/easyctl/pkg/util/strings"
	"gopkg.in/yaml.v2"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
)

const SurviveListFilePath = "server-list.txt"

type Config struct {
	Ping []struct {
		Address string `yaml:"address"`
		Start   int    `yaml:"start"`
		End     int    `yaml:"end"`
		Port    int    `yaml:"port"`
	} `yaml:"ping"`
}

type Manger struct {
	Handler      HandlerInterface
	CheckCount   int
	CheckTimeout time.Duration
}

// Run 执行指令
func Run(item command.OperationItem) command.RunErr {

	serverList, err := ParsePingItems(item.B, item.Logger)
	if err != nil {
		return command.RunErr{Err: err}
	}

	m := Manger{
		Handler:      getHandlerInterface(item.Interface),
		CheckTimeout: time.Second,
		CheckCount:   1,
	}

	return command.RunErr{Err: m.getSurviveList(serverList, item.Logger)}
}

// ParsePingItems 解析探活地址区间
func ParsePingItems(b []byte, logger *logrus.Logger) (strings2.IPS, error) {
	var address strings2.IPS

	item := &Config{}

	err := yaml.Unmarshal(b, item)
	if err != nil {
		return nil, command.RunErr{Err: err}
	}

	for _, v := range item.Ping {
		// 192.168.1. -> 192.168.1
		logger.Infof("%s 合法性检测...", v.Address)
		if strings.HasSuffix(v.Address, ".") {
			v.Address = strings.TrimSuffix(v.Address, ".")
		}

		if ip := net.ParseIP(fmt.Sprintf("%s.1", v.Address)); ip == nil {
			logger.Errorf("[interrupt] %s 格式非法", v.Address)
			continue
		}

		logger.Infof("[pass] %s 合法性检测通过!", v.Address)

		logger.Info("ip地址区间合法性检测...")
		if v.Start > 255 || v.Start < 1 || v.End > 255 || v.End < 1 {
			logger.Errorf("[interrupt] ip地址区间不合法，start、end值有效取值区间为[1-255]!")
			continue
		}

		if v.Start > v.End {
			logger.Errorf("[interrupt] ip地址区间不合法，start值不应大于end值!")
			continue
		}

		logger.Info("[pass] ip地址区间合法性检测通过!")

		for i := v.Start; i <= v.End; i++ {
			address = append(address, fmt.Sprintf("%s.%d:%d", v.Address, i, v.Port))
		}
	}

	sort.Sort(address)

	return address, nil
}

func (m Manger) getSurviveList(addresses []string, logger *logrus.Logger) error {

	logger.Infof("生成活体列表文件: %s", SurviveListFilePath)

	surviveList := strings2.IPS{}
	surviveCh := make(chan string, len(addresses))

	logger = log.SetDefault(logger)

	wg := &sync.WaitGroup{}
	wg.Add(len(addresses))

	for _, v := range addresses {
		go func(address string) {
			var telnetErr error
			// 192.168.1.1:22 -> []string{"192.168.1.1", "22"}
			slice := strings.Split(address, ":")
			ip := slice[0]
			logger.Debugf("ping [%s]", ip)

			// ICMP探测
			pingErr := m.Handler.ICMP(ip, m.CheckCount, m.CheckTimeout, true)
			if pingErr != nil {
				logger.Debugf("[%s] ping err: %s", ip, pingErr)
			}

			// 端口可达探测
			if len(slice) > 1 && slice[1] != "0" {
				logger.Debugf("探测端口监听: %s", address)
				telnetErr = m.Handler.Telnet("tcp", address, m.CheckTimeout)
			}

			// 1.ping不通 telnet通
			// 2.ping通telnet通
			// 3.ping通，不进行telnet探测
			case1 := pingErr != nil && telnetErr == nil && len(slice) > 1 && slice[1] != "0"
			if case1 || (pingErr == nil && telnetErr == nil) {
				surviveCh <- ip
			}
			defer wg.Done()
		}(v)
	}

	wg.Wait()
	close(surviveCh)

	for v := range surviveCh {
		surviveList = append(surviveList, v)
	}

	sort.Sort(surviveList)
	logger.Info("探活列表获取完毕")
	return file.WriteWithIPS(surviveList, SurviveListFilePath)
}

func getHandlerInterface(i interface{}) HandlerInterface {
	handlerInterface, _ := i.(HandlerInterface)
	if handlerInterface == nil {
		return new(Handler)
	}
	return handlerInterface
}

package exec

import (
	"errors"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/file"
	strings2 "github.com/weiliang-ms/easyctl/pkg/util/strings"
	"gopkg.in/yaml.v2"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
)

const SurviveListFilePath = "server-list.txt"

type PingItem struct {
	Ping []struct {
		Address string `yaml:"address"`
		Start   int    `yaml:"start"`
		End     int    `yaml:"end"`
		Port    int    `yaml:"port"`
	} `yaml:"ping"`
}

// Ping 执行指令
func Ping(item command.OperationItem) command.RunErr {

	serverList, err := ParsePingItems(item.B, item.Logger)
	if err != nil {
		return command.RunErr{Err: err}
	}

	surviveList, err := getSurviveList(serverList, item.Logger)
	if err != nil {
		return command.RunErr{Err: err}
	}

	item.Logger.Infof("生成活体列表文件: %s", SurviveListFilePath)
	return command.RunErr{Err: file.WriteWithIPS(surviveList, SurviveListFilePath)}
}

// ParsePingItems 解析探活地址区间
func ParsePingItems(b []byte, logger *logrus.Logger) (strings2.IPS, error) {
	var address strings2.IPS

	item := &PingItem{}

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
			break
		}

		logger.Infof("[pass] %s 合法性检测通过!", v.Address)

		logger.Info("ip地址区间合法性检测...")
		if v.Start > 255 || v.Start < 1 || v.End > 255 || v.End < 1 {
			logger.Errorf("[interrupt] ip地址区间不合法，start、end值有效取值区间为[1-255]!")
			break
		}

		if v.Start > v.End {
			logger.Errorf("[interrupt] ip地址区间不合法，start值不应大于end值!")
			break
		}

		logger.Info("[pass] ip地址区间合法性检测通过!")

		for i := v.Start; i <= v.End; i++ {
			address = append(address, fmt.Sprintf("%s.%d:%d", v.Address, i, v.Port))
		}
	}

	sort.Sort(address)

	return address, nil
}

func getSurviveList(addresses []string, logger *logrus.Logger) (strings2.IPS, error) {

	surviveList := strings2.IPS{}
	surviveCh := make(chan string, len(addresses))

	wg := &sync.WaitGroup{}
	wg.Add(len(addresses))

	for _, v := range addresses {
		go func(address string) {
			var err error
			// 192.168.1.1:22 -> []string{"192.168.1.1", "22"}
			slice := strings.Split(address, ":")
			ip := slice[0]
			logger.Debugf("ping [%s]", ip)

			// icmp探测
			pingErr := icmp(ip, true)
			if pingErr != nil {
				logger.Debugf("[%s] ping err: %s", ip, pingErr)
			}

			// 端口可达探测
			if len(slice) > 1 && slice[1] != "0" {
				logger.Debugf("探测端口监听: %s", address)
				_, err = net.DialTimeout("tcp", address, 1*time.Second)
			}

			// 1.ping不通 telnet通
			// 2.ping通telnet通
			// 3.ping通，不进行telnet探测
			if (pingErr != nil && err == nil) || (pingErr == nil && err == nil) {
				surviveCh <- ip
			}
			defer wg.Done()
		}(v)
	}

	for {
		v, err := readWithSelect(surviveCh)
		if err != nil {
			break
		} else {
			surviveList = append(surviveList, v)
		}
	}

	wg.Wait()

	sort.Sort(surviveList)
	logger.Info("探活列表获取完毕")
	return surviveList, nil
}

func readWithSelect(ch chan string) (x string, err error) {
	timeout := time.NewTimer(time.Second * 2)

	select {
	case x = <-ch:
		return x, nil
	case <-timeout.C:
		return "", errors.New("read time out")
	}
}

func icmp(ip string, privileged bool) error {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		panic(err)
	}
	pinger.Timeout = 1 * time.Second
	pinger.Count = 1
	pinger.SetPrivileged(privileged)
	if err := pinger.Run(); err != nil {
		return err
	}

	if pinger.Statistics().PacketsRecv == 0 {
		return fmt.Errorf("%s无法ping通", ip)
	}

	return nil
}

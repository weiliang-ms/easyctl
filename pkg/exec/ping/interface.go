package ping

import (
	"fmt"
	"github.com/go-ping/ping"
	"github.com/weiliang-ms/easyctl/pkg/util/errors"
	"net"
	"time"
)

//go:generate mockery --name=HandlerInterface
type HandlerInterface interface {
	ICMP(ip string, count int, timeout time.Duration, privileged bool) error
	Telnet(protocol string, address string, timeout time.Duration) error
}

type Handler struct{}

func (h Handler) ICMP(ip string, count int, timeout time.Duration, privileged bool) error {
	p, err := ping.NewPinger(ip)
	if err != nil {
		return err
	}
	p.Timeout = timeout
	p.Count = count
	p.SetPrivileged(privileged)
	// todo: this case
	_ = p.Run()
	//if err := p.Run(); err != nil {
	//	return err
	//}

	return errors.FalseConditionErr(p.Statistics().PacketsRecv == 0, fmt.Sprintf("%s无法ping通", ip))
}

func (h Handler) Telnet(protocol string, address string, timeout time.Duration) error {
	_, err := net.DialTimeout(protocol, address, timeout)
	return err
}

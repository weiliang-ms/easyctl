package util

import (
	"bufio"
	"context"
	"fmt"
	"github.com/yahoo/vssh"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type SSHObject struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (ssh SSHObject) ExecuteOriginCmd(cmd string) {
	vs := vssh.New().Start()
	config := vssh.GetConfigUserPass(ssh.Username, ssh.Password)
	//for _, addr := range []string{"192.168.239.131:22", "192.168.239.141:22"} {
	//	vs.AddClient(addr, config, vssh.SetMaxSessions(2))
	//}
	vs.AddClient(fmt.Sprintf(ssh.Host+":"+ssh.Port), config, vssh.SetMaxSessions(2))
	vs.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd = "ping -c 4 192.168.239.2"
	timeout, _ := time.ParseDuration("6s")
	respChan := vs.Run(ctx, cmd, timeout)

	for resp := range respChan {
		if err := resp.Err(); err != nil {
			log.Println(err)
			continue
		}

		outTxt, errTxt, _ := resp.GetText(vs)
		fmt.Println(outTxt, errTxt, resp.ExitStatus())
	}
}
func ReadSSHInfoFromFile(hostsFile string) (objects []SSHObject) {
	f, err := os.OpenFile(hostsFile, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		for _, v := range strings.Split(line, " ") {
			fmt.Printf("%v", v)
			var object SSHObject
			host := strings.Contains(v, "host")
			port := strings.Contains(v, "port")
			username := strings.Contains(v, "username")
			password := strings.Contains(v, "password")
			if host && port && password && username {
				object.Host = strings.Trim(line, "host=")
				object.Port = strings.Trim(line, "port=")
				object.Username = strings.Trim(line, "username=")
				object.Password = strings.Trim(line, "password=")
				objects = append(objects, object)
			}

		}
		if err != nil || io.EOF == err {
			break
		}
	}

	fmt.Printf("+%v", objects)
	return objects
}

///*
//Copyright 2020 The KubeSphere Authors.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//*/
//
package ssh

//
//import (
//	"fmt"
//	"golang.org/x/crypto/ssh"
//	"net"
//	"strconv"
//	"sync"
//	"time"
//)
//
//type Server struct {
//	ID             int    `json:"-"`
//	Host           string `yaml:"host"`
//	Port           string `yaml:"port"`
//	Username       string `yaml:"username"`
//	Password       string `yaml:"password"`
//	PrivateKeyPath string `yaml:"privateKeyPath,omitempty" json:"privateKeyPath,omitempty"`
//}
//
//type Servers struct {
//	Server []Server `yaml:"server"`
//}
//
//type Dialer struct {
//	lock        sync.Mutex
//	connections map[int]Connection
//}
//
//func NewDialer() *Dialer {
//	return &Dialer{
//		connections: make(map[int]Connection),
//	}
//}
//
//func (dialer *Dialer) Connect(server Server) (Connection, error) {
//	var err error
//
//	dialer.lock.Lock()
//	defer dialer.lock.Unlock()
//
//	conn, _ := dialer.connections[server.ID]
//	port, _ := strconv.Atoi(server.Port)
//	opts := Cfg{
//		Username: server.Username,
//		Port:     port,
//		Address:  server.Host,
//		Password: server.Password,
//		KeyFile:  server.PrivateKeyPath,
//		Timeout:  30 * time.Second,
//	}
//	conn, err = NewConnection(opts)
//	if err != nil {
//		return nil, err
//	}
//	dialer.connections[server.ID] = conn
//
//	return conn, nil
//}
//
//func (server *Server) SSHConnect() (*ssh.Session, error) {
//	var (
//		auth         []ssh.AuthMethod
//		addr         string
//		clientConfig *ssh.ClientConfig
//		client       *ssh.Client
//		session      *ssh.Session
//		err          error
//	)
//	// get auth method
//	auth = make([]ssh.AuthMethod, 0)
//	auth = append(auth, ssh.Password(server.Password))
//
//	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
//		return nil
//	}
//
//	clientConfig = &ssh.ClientConfig{
//		User: server.Username,
//		Auth: auth,
//		// Timeout:             30 * time.Second,
//		HostKeyCallback: hostKeyCallbk,
//	}
//
//	// connet to ssh
//	addr = fmt.Sprintf("%s:%s", server.Host, server.Port)
//
//	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
//		return nil, err
//	}
//
//	// create session
//	if session, err = client.NewSession(); err != nil {
//		return nil, err
//	}
//
//	return session, nil
//}

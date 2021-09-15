package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
)

type Server struct {
	Host           string
	Port           string
	Username       string
	Password       string
	PrivateKeyPath string
}

func (server Server) sshConnect() (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	if server.Password != "" {
		auth = append(auth, ssh.Password(server.Password))
	}

	//auth = append(auth, ssh.PublicKeys(server.PrivateKeyPath))
	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User: server.Username,
		Auth: auth,
		// Timeout:             30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%s", server.Host, server.Port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

//func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
//
//	if err != nil {
//		log.Fatal("find key's home dir failed", err)
//	}
//	key, err := ioutil.ReadFile(keyPath)
//	if err != nil {
//		log.Fatal("ssh key file read failed", err)
//	}
//	// Create the Signer for this private key.
//	signer, err := ssh.ParsePrivateKey(key)
//	if err != nil {
//		log.Fatal("ssh key signer failed", err)
//	}
//	return ssh.PublicKeys(signer)
//}

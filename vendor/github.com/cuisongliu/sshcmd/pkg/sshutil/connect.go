package sshutil

import (
	"fmt"
	"github.com/wonderivan/logger"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

/**
这里主要是做连接ssh操作的
*/
func (ss *SSH) Connect(host string) (*ssh.Session, error) {
	auth := []ssh.AuthMethod{ss.sshAuthMethod(ss.Password, ss.PkFile)}
	config := ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}

	clientConfig := &ssh.ClientConfig{
		User:    ss.User,
		Auth:    auth,
		Timeout: time.Duration(1) * time.Minute,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := ss.addrReformat(host)
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}

	return session, nil
}

func (ss *SSH) sshAuthMethod(passwd, pkFile string) ssh.AuthMethod {
	var am ssh.AuthMethod
	if passwd != "" {
		am = ssh.Password(passwd)
	} else {
		pkData := ss.readFile(pkFile)
		pk, _ := ssh.ParsePrivateKey([]byte(pkData))
		am = ssh.PublicKeys(pk)
	}
	return am
}

func (ss *SSH) readFile(name string) string {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		logger.Error("[globals] read file err is : %s", err)
		return ""
	}
	return string(content)
}

func (ss *SSH) addrReformat(host string) string {
	if strings.Index(host, ":") == -1 {
		host = fmt.Sprintf("%s:22", host)
	}
	return host
}

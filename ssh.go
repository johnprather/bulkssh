package bulkssh

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func sshAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func newSSHClientConfig(user *string, pass *string, agent *bool) *ssh.ClientConfig {
	sshConfig := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{},
	}
	if agent != nil && *agent {
		sshConfig.Auth = append(sshConfig.Auth, sshAgent())
	}
	if pass != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.Password(*pass))
	}
	return sshConfig
}

func sshInit(req *Request) (*ssh.Client, error) {
	sshConfig := newSSHClientConfig(req.User, req.Password, req.Agent)
	hostPort := fmt.Sprintf("%s:%d", req.Hostname, req.Port)
	var err error
	var conn net.Conn
	if req.ConnectTimeout > 0 {
		conn, err = net.DialTimeout("tcp", hostPort,
			time.Duration(req.ConnectTimeout)*time.Second)
	} else {
		conn, err = net.Dial("tcp", hostPort)
	}
	if err != nil {
		return nil, err
	}
	connection, chans, reqs, err := ssh.NewClientConn(conn, hostPort, sshConfig)
	//connection, err := ssh.Dial("tcp", hostPort, sshConfig)
	if err != nil {
		return nil, err
	}
	return ssh.NewClient(connection, chans, reqs), nil
}

func sshRun(connection *ssh.Client, command *Command, timeout int) (*string, error) {
	session, err := connection.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	//fmt.Printf("Established session with %s\n", m.hostname)
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err = session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}

	resStr := make(chan bool)
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	go func() {
		session.Run(command.Command)
		resStr <- true
	}()
	for {
		select {
		case b := <-resStr:
			if b {
				someout := stdoutBuf.String()
				return &someout, nil
			}
		case <-time.After(time.Duration(timeout) * time.Second):
			if timeout > 0 {
				err = fmt.Errorf("Command timed out after %d seconds", timeout)
				return nil, err
			}
		}
	}
}

func sshDisconnect(connection *ssh.Client) {
	connection.Close()
	connection.Conn.Close()
}

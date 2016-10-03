package bulkssh

import (
	"bytes"
	"fmt"
	"net"
	"os"

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
	connection, err := ssh.Dial("tcp", hostPort, sshConfig)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func sshRun(connection *ssh.Client, command *Command) (*string, error) {
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
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(command.Command)
	someout := stdoutBuf.String()
	return &someout, nil
}

func sshDisconnect(connection *ssh.Client) {
	connection.Close()
	connection.Conn.Close()
}

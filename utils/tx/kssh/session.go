package kssh

import (
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type SSHSession struct {
	In      io.Reader
	Out     io.WriteCloser
	Session *ssh.Session
}

// Create a new SSH session (Dial, open in/out pipes and start the shell)
// pass the authntication details in sshConfig
func NewSSHSession(host string, sshConfig *ssh.ClientConfig) (*SSHSession, error) {
	if !strings.Contains(host, ":") {
		return nil, fmt.Errorf("include the port in the host: %s", host)
	}

	connection, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %s", err)
	}
	session, err := connection.NewSession()
	if err != nil {
		return nil, err
	}
	sshIn, err := session.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("session stdout: %s", err)
	}
	sshOut, err := session.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("session stdin: %s", err)
	}
	// sshIn2, err := session.StderrPipe()
	// if err != nil {
	// 	return nil, fmt.Errorf("session stderr: %s", err)
	// }
	// Request PTY (required for srl)
	modes := ssh.TerminalModes{
		ssh.ECHO: 1, // disable echo
	}
	err = session.RequestPty("dumb", 24, 100, modes)
	if err != nil {
		err2 := session.Close()
		return nil, fmt.Errorf("pty request failed: %s - %s", err, err2)
	}

	if err := session.Shell(); err != nil {
		err2 := session.Close()
		return nil, fmt.Errorf("session shell: %s - %s", err, err2)
	}

	return &SSHSession{
		Session: session,
		In:      sshIn,
		Out:     sshOut,
	}, nil
}

func (ses *SSHSession) Writeln(command string) (int, error) {
	return ses.Out.Write([]byte(command + "\r"))
}

func (ses *SSHSession) Close() error {
	log.Debugf("Closing session")
	return ses.Session.Close()
}

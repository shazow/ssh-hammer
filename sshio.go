package main

import (
	"io"

	"golang.org/x/crypto/ssh"
)

type SessionIO struct {
	io.Reader
	io.Writer
}

func NewSessionIO(session *ssh.Session) (r io.Reader, w io.Writer, err error) {
	w, err = session.StdinPipe()
	if err != nil {
		return
	}

	r, err = session.StdoutPipe()
	if err != nil {
		return
	}

	err = session.RequestPty("xterm", 80, 40, ssh.TerminalModes{})
	if err != nil {
		return
	}

	err = session.Shell()
	if err != nil {
		return
	}

	return
}

package main

import (
	"io"

	"golang.org/x/crypto/ssh"
)

type SessionIO struct {
	io.Reader
	io.Writer
}

func NewSessionIO(session *ssh.Session) (err error, r io.Reader, w io.Writer) {
	w, err = session.StdinPipe()
	if err != nil {
		return
	}

	r, err = session.StdoutPipe()
	if err != nil {
		return
	}

	err = session.Shell()
	if err != nil {
		return
	}

	err = session.RequestPty("xterm", 80, 40, ssh.TerminalModes{})
	if err != nil {
		return
	}

	return
}

func doHammerThings(session *ssh.Session) error {
	defer session.Close()

	err, r, w := NewSessionIO(session)
	if err != nil {
		return err
	}

	actor := NewActor(r, w)
	go Echo(actor)

	actor.Wait()
	return nil
}

func Hammer(host string) (chan struct{}, error) {
	done := make(chan struct{})

	fail := func(err error) error {
		close(done)
		return err
	}

	key, err := MakeKey()
	if err != nil {
		return done, fail(err)
	}

	config := &ssh.ClientConfig{
		User: "hammer",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return done, fail(err)
	}

	session, err := conn.NewSession()
	if err != nil {
		return done, fail(err)
	}

	doHammerThings(session)

	return done, nil
}

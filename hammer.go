package main

import (
	"bytes"

	"golang.org/x/crypto/ssh"
)

func doHammerThings(session *ssh.Session) {
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run("ls")

}

func Hammer(host string) chan struct{} {
	done := make(chan struct{})

	key, err := MakeKey()
	if err != nil {
		logger.Errorf("Failed to make key: %s", err)
		close(done)
		return done
	}

	config := &ssh.ClientConfig{
		User: "hammer",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	go func() {
		conn, err := ssh.Dial("tcp", host, config)
		if err != nil {
			logger.Errorf("Failed to connect: %s", err)
			close(done)
			return
		}
		session, err := conn.NewSession()
		if err != nil {
			logger.Errorf("Failed to create session: %s", err)
			close(done)
			return
		}
		doHammerThings(session)
	}()

	return done
}

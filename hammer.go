package main

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type Hammer struct {
	host string
	done chan struct{}
	num  int
}

func NewHammer(host string, num int) *Hammer {
	return &Hammer{
		host: host,
		num:  num,
		done: make(chan struct{}, 1),
	}
}

func (h *Hammer) Wait() {
	<-h.done
}

func (h *Hammer) Stop() {
	close(h.done)
}

func (h *Hammer) Start() error {
	key, err := MakeKey()
	if err != nil {
		return err
	}

	auth := []ssh.AuthMethod{
		ssh.PublicKeys(key),
	}

	for i := 1; i <= h.num; i++ {
		name := fmt.Sprintf("hammer%d", i)
		err := h.connect(name, auth)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Hammer) connect(name string, auth []ssh.AuthMethod) error {
	config := &ssh.ClientConfig{
		User:            name,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", h.host, config)
	if err != nil {
		return err
	}

	session, err := conn.NewSession()
	if err != nil {
		return err
	}

	go func() {
		h.Wait()
		session.Close()
	}()

	r, w, err := NewSessionIO(session)
	if err != nil {
		return err
	}

	go func() {
		Spam(NewActor(r, w))
	}()
	return nil
}

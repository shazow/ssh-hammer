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
	for i := 1; i <= h.num; i++ {
		name := fmt.Sprintf("hammer%d", i)
		err := h.connect(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Hammer) connect(name string) error {
	// TODO: Re-use keys.
	logger.Debugf("MakeKey")
	key, err := MakeKey()
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: name,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	conn, err := ssh.Dial("tcp", h.host, config)
	if err != nil {
		return err
	}

	logger.Debugf("NewSession")
	session, err := conn.NewSession()
	if err != nil {
		return err
	}

	go func() {
		h.Wait()
		session.Close()
	}()

	logger.Debugf("NewSessionIO")
	r, w, err := NewSessionIO(session)
	if err != nil {
		return err
	}

	logger.Debugf("NewActor")
	go func() {
		Spam(NewActor(r, w))
	}()
	return nil
}

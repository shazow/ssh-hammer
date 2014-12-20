package main

import (
	"bufio"
	"io"
)

type Actor struct {
	*bufio.Reader
	out  io.Writer
	done chan struct{}
}

func NewActor(r io.Reader, w io.Writer) *Actor {
	a := Actor{
		Reader: bufio.NewReader(r),
		out:    w,
		done:   make(chan struct{}),
	}
	return &a
}

func (a *Actor) Write(p []byte) (n int, err error) {
	return a.out.Write(p)
}

func (a *Actor) Wait() {
	<-a.done
}

func (a *Actor) Close() {
	close(a.done)
}

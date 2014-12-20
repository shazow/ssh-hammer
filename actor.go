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

func Echo(a *Actor) {
	echo := make(chan string, 10)
	consume := make(chan struct{}, 1)
	stop := func(err error) {
		logger.Errorf("Echo stopped: %s", err)
		close(echo)
		a.Close()
	}

	go func() {
		// Reader
		s := bufio.NewScanner(a)
		for s.Scan() {
			msg := s.Text()
			logger.Debugf("Read: %s", msg)
			echo <- msg
		}

		err := s.Err()
		if err != nil {
			stop(err)
			return
		}
	}()

	go func() {
		// Consume until we stop receiving content.
		defer close(consume)
		for {
			select {
			case <-echo:
			default:
				return
			}
		}
	}()

	go func() {
		// Writer
		<-consume
		logger.Debugf("Starting to echo.")

		for msg := range echo {
			_, err := a.Write([]byte(msg + "\r\n"))
			if err != nil {
				stop(err)
				return
			}
		}

	}()
}

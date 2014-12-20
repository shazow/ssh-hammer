package main

import "time"

func Spam(a *Actor) {
	logger.Infof("Starting to spam.")

	for {
		<-time.After(1 * time.Second)
		_, err := a.Write([]byte("Spam!\r\n"))
		if err != nil {
			a.Close()
			return
		}
	}
}

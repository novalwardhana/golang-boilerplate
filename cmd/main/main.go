package main

import (
	httpHandler "github.com/novalwardhana/golang-boilerplate/cmd/http-handler"
)

func main() {

	start := make(chan int)
	go func() {
		go httpHandler.RunHTTPHandler()
	}()
	<-start

}

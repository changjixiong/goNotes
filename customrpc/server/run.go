package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type appSvr struct {
	servers []Server
	stop    chan struct{}
	destory func()
}

func NewApp() *appSvr {
	a := new(appSvr)
	a.stop = make(chan struct{}, 1)
	return a
}

func (a *appSvr) Run() {
	defer func() {
		if err := recover(); err != nil {

			log.Println("Recover from panic.", err)
		}

		for i := 0; i < len(a.servers); i++ {

			log.Println("close server", i)
			a.servers[i].Close()
		}
	}()

	for i := 0; i < len(a.servers); i++ {

		log.Println("start server", i)
		a.servers[i].Start()
	}

	log.Println(len(a.servers), " services started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

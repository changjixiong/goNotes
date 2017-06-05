package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"net/http"

	consulapi "github.com/hashicorp/consul/api"
)

const RECV_BUF_LEN = 1024

func consulCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "consulCheck")
}

func registerServer() {

	config := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(config)

	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	checkPort := 8080

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "serverNode_2"
	registration.Name = "serverNode"
	registration.Port = 9527
	registration.Tags = []string{"serverNode"}
	registration.Address = "127.0.0.1"
	registration.Check = &consulapi.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, checkPort, "/check"),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务
	}

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		log.Fatal("register server error : ", err)
	}

	http.HandleFunc("/check", consulCheck)
	http.ListenAndServe(fmt.Sprintf(":%d", checkPort), nil)

}

func main() {

	go registerServer()

	ln, err := net.Listen("tcp", "0.0.0.0:9527")

	if nil != err {
		panic("Error: " + err.Error())
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			panic("Error: " + err.Error())
		}

		go EchoServer(conn)
	}

}

func EchoServer(conn net.Conn) {
	buf := make([]byte, RECV_BUF_LEN)
	defer conn.Close()

	for {
		n, err := conn.Read(buf)
		switch err {
		case nil:
			log.Println("get and echo:", "EchoServer "+string(buf[0:n]))
			conn.Write(append([]byte("EchoServer "), buf[0:n]...))
		case io.EOF:
			log.Printf("Warning: End of data: %s\n", err)
			return
		default:
			log.Printf("Error: Reading data: %s\n", err)
			return
		}
	}
}

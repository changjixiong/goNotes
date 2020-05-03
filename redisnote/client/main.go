package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

const RECV_BUF_LEN = 1024

func main() {

	client, err := consulapi.NewClient(consulapi.DefaultConfig())

	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	for {

		time.Sleep(time.Second * 3)
		var services map[string]*consulapi.AgentService
		var err error

		services, err = client.Agent().Services()

		log.Println("services", strings.Repeat("-", 80))
		for _, service := range services {
			log.Println(service)
		}

		if nil != err {
			log.Println("in consual list Services:", err)
			continue
		}

		if _, found := services["rankNode_1"]; !found {
			log.Println("rankNode_1 not found")
			continue
		}
		log.Println("choose", strings.Repeat("-", 80))
		log.Println("rankNode_1", services["rankNode_1"])
		sendData(services["rankNode_1"])

	}
}

func sendData(service *consulapi.AgentService) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", service.Address, service.Port))

	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	buf := make([]byte, RECV_BUF_LEN)

	r := rand.New(rand.NewSource(time.Now().Unix()))

	for {
		msg := ""
		switch r.Intn(2) {
		case 0:
			msg = fmt.Sprintf(`{"func_name":"GetPlayerByLvRank","params":[%d,%d]}`, 0, 3)
		case 1:
			msg = fmt.Sprintf(`{"func_name":"AddPlayerExp","params":[%d,%d]}`, 1+r.Intn(5), r.Intn(100))
		}

		n, err := conn.Write([]byte(msg))
		if err != nil {
			println("Write Buffer Error:", err.Error())
			break
		}

		n, err = conn.Read(buf)
		if err != nil {
			println("Read Buffer Error:", err.Error())
			break
		}
		log.Println("get:", string(buf[0:n]))

		//等一秒钟
		time.Sleep(time.Second)
	}
}

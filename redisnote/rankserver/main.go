package main

import (
	"fmt"
	"goNotes/redisnote/playercache"
	"goNotes/redisnote/rankservice"
	"goNotes/reflectinvoke"
	"io"
	"log"
	"net"
	"net/http"

	consulapi "github.com/hashicorp/consul/api"
	redis "gopkg.in/redis.v4"
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

	checkPort := 8081

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "rankNode_1"
	registration.Name = "rankNode"
	registration.Port = 9528
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

func RankServer(conn net.Conn) {
	buf := make([]byte, RECV_BUF_LEN)
	defer conn.Close()

	for {
		n, err := conn.Read(buf)
		switch err {
		case nil:
			log.Println("get:", string(buf[0:n]))

			resultdata := reflectinvoke.InvokeByJson([]byte(buf[0:n]))
			conn.Write(resultdata)
		case io.EOF:
			log.Printf("Warning: End of data: %s\n", err)
			return
		default:
			log.Printf("Error: Reading data: %s\n", err)
			return
		}
	}
}

func main() {

	rdsClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password: "123456",
		DB:       0,
	})
	playercache.Rds = rdsClient
	rankservice.Rds = rdsClient

	for i := 0; i < 5; i++ {
		player := &playercache.PlayerInfo{
			PlayerID:   i + 1,
			PlayerName: fmt.Sprintf("玩家%d", i+1),
			Lv:         1,
			Online:     true,
		}

		playercache.SetPlayerInfo(player)

	}

	// 用json字符串调用
	reflectinvoke.RegisterMethod(rankservice.DefaultRankService)

	go registerServer()

	ln, err := net.Listen("tcp", "0.0.0.0:9528")

	if nil != err {
		panic("Error: " + err.Error())
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			panic("Error: " + err.Error())
		}

		go RankServer(conn)
	}

}

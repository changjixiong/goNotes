package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

import ini "gopkg.in/ini.v1"

var cfg *ini.File

func main() {
	var err error
	cfg, err = ini.Load("debuggerInTheMiddle.ini")
	listener, err := net.Listen("tcp", cfg.Section("middle").Key("address").String())
	if err != nil {
		panic("error listening:" + err.Error())
	}
	fmt.Println("Starting the server")

	for {
		conn, err := listener.Accept() //接受连接
		if err != nil {
			panic("Error accept:" + err.Error())
		}
		fmt.Println("")
		fmt.Println("")
		fmt.Println("---------- Accepted the Connection :", conn.RemoteAddr(), " ----------")
		go MiddleServer(conn)
	}
}

func MiddleServer(conn net.Conn) {
	buf := make([]byte, 1024)
	defer conn.Close()

	for {
		n, err := conn.Read(buf)
		switch err {
		case nil:
			sendToRemote(conn, buf[0:n])
		case io.EOF:
			fmt.Printf("########Warning: End of data: %s \n", err)
			fmt.Println("########client may by shut down")
			return
		default:
			fmt.Printf("########Error: Reading data : %s \n", err)
			return
		}
	}
}

func sendToRemote(client net.Conn, message []byte) {
	conn, err := net.Dial("tcp", cfg.Section("remote").Key("address").String())
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	len := len(message)
	haveSend := 0

	for haveSend < len {
		n, err := conn.Write(message[haveSend:])
		if err != nil {
			println("Write Buffer Error:", err.Error())
			return
		}

		haveSend += n
	}

	fmt.Println("sendToRemote ----------------->")
	fmt.Println(string(message))

	//从服务器端收字符串
	n, err := conn.Read(buf)
	if err != nil {
		println("Read Buffer Error:", err.Error())
		return
	}

	fmt.Println("<------------------ getToRemote")
	fmt.Println(string(buf))
	fmt.Println(strings.Repeat("-", 64))
	client.Write(buf[0:n])
}

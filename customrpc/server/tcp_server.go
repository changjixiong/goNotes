package main

import (
	"io"
	"log"
	"net"
)

const RECV_BUF_LEN = 1024



type TCPServer struct {
	addr string
}

func (t *TCPServer) Start() {

	ln, err := net.Listen("tcp", t.addr)

	if nil != err {
		panic("Error: " + err.Error())
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			panic("Error: " + err.Error())
		}

		go t.run(conn)
	}
}

func (t *TCPServer) run(conn net.Conn) {
	buf := make([]byte, RECV_BUF_LEN)
	defer conn.Close()

	for {
		n, err := conn.Read(buf)
		switch err {
		case nil:
			log.Println("get:", string(buf[0:n]))

			resultdata := reflectinvoker.InvokeByJson([]byte(buf[0:n]))
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

func (t *TCPServer) Close() {

}

func NewTCPServer(addr string) *TCPServer {

	s := &TCPServer{
		addr: addr,
	}

	return s

}

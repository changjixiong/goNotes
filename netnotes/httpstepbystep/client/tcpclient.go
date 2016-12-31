package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("net.Dial 127.0.0.1:8080 :", err)
		os.Exit(1)
	}

	for {
		fmt.Fprintf(conn, "GET / HTTP/1.1\nHost: 127.0.0.1:8080\nConnection: keep-alive\n\n")
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)

		switch err {
		case nil:
		case io.EOF:
			fmt.Println("detected closed LAN connection")
			return
		default:
			fmt.Println("conn.Read", err)
			return
		}

		fmt.Printf("recv %d:%s", n, string(buf))

		time.Sleep(time.Second * 2000)
	}

}

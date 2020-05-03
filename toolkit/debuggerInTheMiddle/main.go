package main

import (
	"fmt"
	"io"
	"net"

	ini "gopkg.in/ini.v1"
)

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

		fmt.Println("---------- Accepted the Connection :", conn.RemoteAddr(), " ----------")
		remoteConn, err := net.Dial("tcp", cfg.Section("remote").Key("address").String())
		if err != nil {
			panic(err.Error())
		}

		go MiddleServer(conn, remoteConn)
	}
}

func MiddleServer(clientConn net.Conn, remoteConn net.Conn) {

	go transfer(clientConn, remoteConn)
	go transfer(remoteConn, clientConn)

}

func sendAll(conn net.Conn, msg []byte) error {
	msglen := len(msg)
	haveSend := 0

	for haveSend < msglen {
		n, err := conn.Write(msg[haveSend:])
		if err != nil {

			return err
		}

		haveSend += n
	}

	return nil
}

func transfer(srcConn, destConn net.Conn) {

	defer func() {
		srcConn.Close()
		destConn.Close()
	}()

	const (
		sendout = 1
		getback = 2
	)

	transferDir := 0

	if srcConn.LocalAddr().String() == cfg.Section("middle").Key("address").String() {
		transferDir = sendout
	} else {
		transferDir = getback
	}

	for {

		buf := make([]byte, 4096)

		n, err := srcConn.Read(buf)

		switch err {
		case nil:
		case io.EOF:
			fmt.Printf("########Warning: End of data: %s \n", err)
			fmt.Println("########client may by shut down")
			return
		default:
			fmt.Printf("########Error: Reading data : %s \n", err)
			return
		}

		if "1" == cfg.Section("debug").Key("printclientsend").String() &&
			sendout == transferDir {
			fmt.Println("----  getFromClient  -------------->")
			fmt.Println(string(buf[0:n]))
		}

		if "1" == cfg.Section("debug").Key("printserversend").String() &&
			getback == transferDir {
			fmt.Println("<--------------  getFromServer  ----")
			fmt.Println(string(buf[0:n]))
		}

		err = sendAll(destConn, buf[0:n])

		if nil != err {
			println(err)
			return
		}
	}

}

package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

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
		fmt.Println("")
		fmt.Println("")
		fmt.Println("---------- Accepted the Connection :", conn.RemoteAddr(), " ----------")
		go MiddleServer(conn)
	}
}

func MiddleServer(clientConn net.Conn) {
	buf := make([]byte, 8192)
	defer clientConn.Close()

	remoteConn, err := net.Dial("tcp", cfg.Section("remote").Key("address").String())
	if err != nil {
		panic(err.Error())
	}
	defer remoteConn.Close()

	clientOutChan := make(chan []byte, 128)
	clientBackChan := make(chan []byte, 64)

	for {
		n, err := clientConn.Read(buf)
		go transfer(clientConn, remoteConn, clientOutChan, clientBackChan)
		switch err {
		case nil:
			// sendToRemote(clientConn, remoteConn, buf[0:n])
			clientOutChan <- buf[0:n]
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

func transfer(client net.Conn, remote net.Conn,
	clientOutChan chan []byte,
	clientBackChan chan []byte) {

	ticker := time.NewTicker(time.Microsecond * 1)

	defer func() {
		ticker.Stop()
	}()

	for {

		select {
		case msg, ok := <-clientOutChan:

			if !ok {
				return
			}
			msglen := len(msg)
			haveSend := 0

			for haveSend < msglen {
				n, err := remote.Write(msg[haveSend:])
				if err != nil {
					println("Write Buffer Error:", err.Error())
					return
				}

				haveSend += n
			}
			fmt.Println("sendToRemote ----------------->")
			fmt.Println(string(msg))
		case <-ticker.C:
		}

		buf := make([]byte, 4096)
		//从服务器端收字符串
		n, err := remote.Read(buf)
		if err != nil {
			println("Read Buffer Error:", err.Error())
			return
		}

		fmt.Println("<------------------ getToRemote")
		fmt.Println(string(buf))
		// if nil != file {
		// 	file.WriteString("<------------------ getToRemote\n")
		// 	file.Write(buf[0:n])
		// }
		fmt.Println(strings.Repeat("-", 64))
		client.Write(buf[0:n])
	}

}

func sendToRemote(client net.Conn, remote net.Conn, message []byte) {

	msglen := len(message)
	haveSend := 0

	for haveSend < msglen {
		n, err := remote.Write(message[haveSend:])
		if err != nil {
			println("Write Buffer Error:", err.Error())
			return
		}

		haveSend += n
	}

	bnr := bufio.NewReader(strings.NewReader(string(message)))

	contentLen := -1
	nextIsBody := false
	n := 0
	for {
		line, err := bnr.ReadString('\n')
		n += 1
		fmt.Println(n)
		if nil != err {
			fmt.Println(err)
			break
		} else {
			if strings.Contains(line, "Content-Length") {
				contentLen = getContentLength(line)
				fmt.Println("contentLen:", contentLen)
			}

			if !nextIsBody && line == "\r\n" {
				nextIsBody = true
				fmt.Println(line, "nextIsBody")
			}

		}
	}

	buf := make([]byte, 1024)

	fileName := ""

	var file *os.File

	if "1" == cfg.Section("config").Key("logout").String() {
		now := time.Now()
		fileName = fmt.Sprintf("%s.%d", now.Format("20060102_150405"), now.Nanosecond()/1000000)
		file, err := os.Create(fileName)
		defer file.Close()
		if nil != err {
			fmt.Println(err)
		}
	}

	fmt.Println("sendToRemote ----------------->")
	fmt.Println(string(message))

	if nil != file {
		file.WriteString("sendToRemote ----------------->\n")
		file.Write(message)
	}

	//从服务器端收字符串
	n, err := remote.Read(buf)
	if err != nil {
		println("Read Buffer Error:", err.Error())
		return
	}

	fmt.Println("<------------------ getToRemote")
	fmt.Println(string(buf))
	if nil != file {
		file.WriteString("<------------------ getToRemote\n")
		file.Write(buf[0:n])
	}
	fmt.Println(strings.Repeat("-", 64))
	client.Write(buf[0:n])
}

func getContentLength(str string) int {
	strs := strings.Split(str, " ")
	n, err := strconv.Atoi(strs[1][0 : len(strs[1])-2])
	if nil != err {
		fmt.Println(str, err)
	}
	return n
}

/*

 */

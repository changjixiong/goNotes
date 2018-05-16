package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"time"
)

type Server struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
type conn struct {
	// server is the server on which the connection arrived.
	// Immutable; never nil.
	server *Server

	// rwc is the underlying network connection.
	// This is never wrapped by other types and is the value given out
	// to CloseNotifier callers. It is usually of type *net.TCPConn or
	// *tls.Conn.
	rwc net.Conn
}

func (c *conn) serve() {

	defer c.rwc.Close()
	buf := make([]byte, 1024)

	for {
		n, err := c.rwc.Read(buf)

		if nil != err {
			fmt.Println("conn.Read err:", err)
			return
		}

		fmt.Println("recv:\n", string(buf[0:n]))
		// data := "HTTP/1.1 200 OK\n" +
		// 	"Date: Wed, 14 Dec 2016 09:50:53 GMT\n" +
		// 	"Content-Length: 4\n" +
		// 	"Content-Type: text/plain; charset=utf-8\n\n" +
		// 	"abcd\n\n"
		data := genResponse("")
		c.rwc.Write([]byte(data))
		fmt.Println("send", data)
	}

}
func (srv *Server) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: srv,
		rwc:    rwc,
	}

	return c
}

func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	fmt.Println("Serve at :", l.Addr())
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, err := l.Accept()
		if err != nil {
			//这个err 是一个 OpError类型，实现了net.Error接口
			//Temporary 临时错误，可重试

			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				fmt.Printf("http: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0
		// go handleConnection(conn)
		c := srv.newConn(rw)
		go c.serve()
	}
}

func handleConnection(conn net.Conn) {

	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)

		if err == io.EOF {
			fmt.Println("conn.Read io.EOF")
			return // don't reply
		}
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			fmt.Println("conn.Read neterr.Timeout()")
			return // don't reply
		}

		fmt.Println("recv:\n", string(buf[0:n]))
		data := genResponse("")
		conn.Write([]byte(data))
		fmt.Println("send", data)
	}
}

func genResponse(content string) string {

	data := "HTTP/1.1 200 OK\r\n"
	data += "Date: " + time.Now().String() + "\r\n"
	data += fmt.Sprintf("Content-Length: %d\r\n", len(content))
	data += "Content-Type: text/plain; charset=utf-8\r\n"
	data += "\r\n"
	data += content + "\r\n"
	data += "\r\n"

	return data
}

func main() {

	port := 0
	flag.IntVar(&port, "port", 8088, "port")
	server := &Server{Addr: fmt.Sprintf(":%d", port)}
	server.ListenAndServe()

}

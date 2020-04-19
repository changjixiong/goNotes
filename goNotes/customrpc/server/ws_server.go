package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	WriteWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 120 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (pongWait * 9) / 10
)

type WSServer struct {
	addr string

	maxMsgLen   int64
	httpTimeout time.Duration

	ln net.Listener

	upgrader websocket.Upgrader
}

func NewWSServer(addr string) *WSServer {
	server := &WSServer{
		addr:        addr,
		maxMsgLen:   512,
		httpTimeout: 5 * time.Second,
	}

	return server
}

func (server *WSServer) run(conn *websocket.Conn, writeChan chan []byte) {

	// log.Println(conn.RemoteAddr(), "run")

	go func() {
		ticker := time.NewTicker(PingPeriod)
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recover from panic:", err)

			}

		}()

		for {
			select {
			case b, ok := <-writeChan:
				conn.SetWriteDeadline(time.Now().Add(WriteWait))
				if b == nil || !ok {
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				err := conn.WriteMessage(websocket.BinaryMessage, b)
				if err != nil {
					log.Println("write msg error :", err)
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					break
				}

			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(WriteWait))
				if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}
			}
		}
	}()

	for {
		_, data, err := conn.ReadMessage()

		if err != nil {
			log.Println("read message error: ", err, conn.RemoteAddr().String())
			writeChan <- nil
			return
		}

		// fmt.Println("get", string(data))
		resultdata := reflectinvoker.InvokeByJson(data)
		writeChan <- resultdata
	}

}

// new client
func (server *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		http.Error(w, "Method upgrade failed.", 500)
		return
	}

	conn.SetReadLimit(server.maxMsgLen)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	writeChan := make(chan []byte, 512)
	server.run(conn, writeChan)
	conn.Close()
	close(writeChan)

}

func (server *WSServer) Start() {

	ln, err := net.Listen("tcp", server.addr)
	if err != nil {
		log.Fatalln("lister error ", err)
	}

	server.ln = ln

	server.upgrader = websocket.Upgrader{
		HandshakeTimeout: server.httpTimeout,
		CheckOrigin:      func(_ *http.Request) bool { return true },
	}

	httpServer := &http.Server{
		Addr:           server.addr,
		Handler:        server,
		ReadTimeout:    server.httpTimeout,
		WriteTimeout:   server.httpTimeout,
		MaxHeaderBytes: 1024,
	}

	log.Println("Start wsserver at", server.addr)
	go httpServer.Serve(ln)
}

func (server *WSServer) Close() {

}

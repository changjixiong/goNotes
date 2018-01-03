package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
)

type HTTPServer struct {
	Addr    string
	ln      net.Listener
	handler http.Handler
	wg      sync.WaitGroup
}

func NewHTTPServer(addr string) *HTTPServer {
	server := &HTTPServer{
		Addr: addr,
	}

	return server
}

func (server *HTTPServer) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("read data error:", err)
		w.WriteHeader(500)
		w.Write([]byte("server error"))
		return
	}
	if len(requestData) < 1 {
		log.Println("RequestData too short.", string(requestData))

		return
	}

	// log.Println("http in =:", string(requestData))

	resultdata := reflectinvoker.InvokeByJson(requestData)

	w.Write(resultdata)
	// w.Write(requestData)

}

func (server *HTTPServer) Start() {

	log.Println("HTTPServer at", server.Addr)
	ln, err := net.Listen("tcp", server.Addr)
	// ln, err := reuseport.NewReusablePortListener("tcp", server.Addr)
	if err != nil {
		log.Fatalln("Listening failed.", err)

	}

	server.ln = ln
	if server.handler == nil {
		mux := http.NewServeMux()
		mux.HandleFunc("/rpc", server.Serve)
		server.handler = mux
	}

	httpServer := &http.Server{
		Addr:           server.Addr,
		Handler:        server.handler,
		MaxHeaderBytes: 1024,
	}

	go httpServer.Serve(ln)
}

func (server *HTTPServer) Close() {

}

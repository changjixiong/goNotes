package main

import (
	"flag"
	"fmt"
	"goNotes/reflectinvoke"

	"goNotes/customrpc/server/services"
)

var reflectinvoker *reflectinvoke.Reflectinvoker
var port1 = flag.Int("port1", 19001, "port")
var port2 = flag.Int("port2", 19002, "port")
var port3 = flag.Int("port3", 19003, "port")

func main() {
	reflectinvoker = reflectinvoke.NewReflectinvoker()
	// 用json字符串调用
	reflectinvoker.RegisterMethod(services.DefaultHelloService)
	reflectinvoker.RegisterMethod(services.DefaultWorldService)

	app := NewApp()

	app.servers = append(app.servers,
		// NewTCPServer(fmt.Sprintf("0.0.0.0:%d", *port1)),
		NewHTTPServer(fmt.Sprintf("0.0.0.0:%d", *port2)),
		NewWSServer(fmt.Sprintf("0.0.0.0:%d", *port3)))

	app.Run()

}

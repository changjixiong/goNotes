package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const RECV_BUF_LEN = 1024

var wg sync.WaitGroup

var n = flag.Int("n", 200, "Number of requests to run. Default is 200.")
var c = flag.Int("c", 50, "Number of requests to run concurrently. Total number of requests cannot be smaller than the concurrency level. Default is 50.")
var serveraddr = flag.String("serveraddr", "127.0.0.1:19003", "serveraddr")

func main() {

	flag.Parse()

	begin := time.Now()

	sendWSData()

	fmt.Println("cost:", time.Now().Sub(begin))

}

func sendTCPData() {
	conn, err := net.Dial("tcp", *serveraddr)

	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	buf := make([]byte, RECV_BUF_LEN)

	r := rand.New(rand.NewSource(time.Now().Unix()))

	for {
		msg := ""
		switch r.Intn(4) {
		case 0:
			msg = fmt.Sprintf(`{"func_name":"Hello","params":[%d]}`, 1)
		case 1:
			msg = fmt.Sprintf(`{"func_name":"Hello","params":[%d]}`, 2)
		case 2:
			msg = fmt.Sprintf(`{"func_name":"World","params":[%d]}`, 1)
		case 3:
			msg = fmt.Sprintf(`{"func_name":"World","params":[%d]}`, 2)

		}

		n, err := conn.Write([]byte(msg))
		if err != nil {
			println("Write Buffer Error:", err.Error())
			break
		}

		n, err = conn.Read(buf)
		if err != nil {
			println("Read Buffer Error:", err.Error())
			break
		}
		log.Println("get:", string(buf[0:n]))

		//等一秒钟
		time.Sleep(time.Second)
	}
}

func sendWSData() {

	for i := 0; i < *c; i++ {
		wg.Add(1)
		go func(i int) {
			errNum := 0
			okNum := 0
			u := url.URL{Scheme: "ws", Host: *serveraddr, Path: ""}

			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			defer conn.Close()
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}

			msg := `{"func_name":"Hello","params":[1]}`

			for j := 0; j < (*n / (*c)); j++ {
				err = conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
				if err == nil {
					okNum++
				} else {
					errNum++
					break
				}

				_, _, err := conn.ReadMessage()
				if err == nil {
					okNum++
				} else {
					errNum++
					break
				}

			}

			if err := conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
				log.Println("websocket.CloseMessage error")
			}
			if *c <= 50 || i%50 == 0 || errNum > 0 {
				fmt.Println(i, "okNum:", okNum, "errNum:", errNum)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

}

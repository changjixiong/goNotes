package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var done = make(chan struct{})
var host = flag.String("host", "http://127.0.0.1", "服务器地址")
var api = flag.String("api", "foo", "api名称")
var port = flag.String("port", "8000", "端口")
var cnum = flag.Int("cnum", 100, "连接数")

var wg sync.WaitGroup
var mu sync.RWMutex
var muprint sync.RWMutex
var worker = 0

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

type Client struct {
	httpclient *http.Client
	requestNum int
}

func handleSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT,
		syscall.SIGTERM)
	for s := range c {
		fmt.Printf("Signals  get  %s ........... \n", s)
		break
	}
	fmt.Println("close(done)")
	close(done)

}

func MakeURL() string {
	return *host + ":" + *port + "/" + *api
}

func AddWorkerNum() {
	mu.Lock()
	defer mu.Unlock()
	worker += 1

}

func WorkerIsFull() bool {
	mu.Lock()
	defer mu.Unlock()

	return worker >= *cnum
}

func work(httpClient *Client, request *http.Request) {

	defer wg.Done()

	for {

		if cancelled() {
			fmt.Println("work cancel")
			return
		}

		needPrint := false
		if httpClient.requestNum > 200 {
			needPrint = true
			httpClient.requestNum = 0
		} else {
			httpClient.requestNum += 1
		}

		var begin time.Time
		if needPrint {
			begin = time.Now()
		}

		resp, err := httpClient.httpclient.Do(request)
		if needPrint {
			fmt.Println("cost", time.Now().Sub(begin).String())
		}

		if nil != err {
			fmt.Println("work err", err)
		} else {
			//fmt.Println(resp)
			resp.Body.Close()
		}

		time.Sleep(time.Millisecond * 5)
	}

}

func NewClient(url string) {
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept-Encoding", "gzip")
	httpClient := &Client{httpclient: &http.Client{}}

	AddWorkerNum()

	wg.Add(1)
	go work(httpClient, request)

}

func GenMulitClient(num int) {

	fmt.Println("GenMulitClient ...")
	defer func() {
		fmt.Println("GenMulitClient finish")
		wg.Done()
	}()

	for i := 0; i < num; i++ {
		if cancelled() {
			return
		}
		NewClient(MakeURL())
		time.Sleep(time.Millisecond * 100)
	}

	wg.Add(1)
	go CheckClient()

}

func CheckClient() {
	defer func() {
		fmt.Println("CheckClient cancle")
		wg.Done()
	}()

	for {

		if cancelled() {

			break
		} else {
			if WorkerIsFull() {
				continue
			}
			NewClient(MakeURL())
			time.Sleep(time.Millisecond * 20)

		}

	}
}

func main() {
	flag.Parse()
	wg.Add(1)
	go GenMulitClient(*cnum)

	handleSignal()
	wg.Wait()
}

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pb3 "goNotes/grpcnotes/echo"

	"golang.org/x/net/context"
)

const (
	address     = "localhost:%d"
	defaultName = "world"
)

var wg sync.WaitGroup

var n = flag.Int("n", 200, "Number of requests to run. Default is 200.")
var c = flag.Int("c", 50, "Number of requests to run concurrently. Total number of requests cannot be smaller than the concurrency level. Default is 50.")
var serveraddr = flag.String("serveraddr", "127.0.0.1:19000", "serveraddr")

func main() {
	flag.Parse()

	begin := time.Now()
	testgrpcstream()
	fmt.Println("cost:", time.Now().Sub(begin))

}

func testgrpcstream() {
	printErrors := make(chan int, 200)

	for i := 0; i < 200; i++ {
		printErrors <- i
	}

	close(printErrors)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < *c; i++ {

		wg.Add(1)
		go func(i int) {
			t := r.Intn(10) + 7
			time.Sleep(time.Millisecond * time.Duration(t))

			canPrintErrors := true
			sendErrNum := 0
			recvErrNum := 0
			sendokNum := 0
			recvokNum := 0
			emptyNum := 0
			conn, err := grpc.Dial(*serveraddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(keepalive.ClientParameters{}))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}

			rclient := pb3.NewEchoClient(conn)

			stream, err := rclient.SayEchoS(context.Background(), grpc.FailFast(false))

			if nil != err {
				log.Fatalf("SayEchoS: %v", err)
			}

			defer func() {
				stream.CloseSend()
				conn.Close()
			}()

			for j := 0; j < (*n / (*c)); j++ {
				time.Sleep(9)

				if err := stream.Send(&pb3.EchoRequest{Num: 1}); err != nil {

					if err != nil {

						if err == io.EOF {
							emptyNum += 1
							continue
						}

						sendErrNum += 1

						if canPrintErrors {
							if _, ok := <-printErrors; ok {
								log.Println("stream.Send EchoRequest:", err)
							} else {
								canPrintErrors = false
							}
						}

					}
				} else {
					sendokNum += 1
				}

				time.Sleep(11)
				_, err := stream.Recv()

				if err != nil {

					if err == io.EOF {
						emptyNum += 1
						continue
					}

					recvErrNum += 1

					if canPrintErrors {
						if _, ok := <-printErrors; ok {
							log.Println("stream.Recv EchoReply:", err)
						} else {
							canPrintErrors = false
						}
					}

				} else {

					recvokNum += 1
				}

			}
			if *c <= 50 || i%50 == 0 || sendErrNum > 0 || recvErrNum > 0 || emptyNum > 0 {
				fmt.Println(i, "sendokNum:", sendokNum, "recvokNum", recvokNum, "sendErrNum:", sendErrNum, "recvErrNum:", recvErrNum, "emptyNum:", emptyNum)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
}

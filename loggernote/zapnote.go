package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/uber-go/zap"
)

func handleSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT,
		//syscall.SIGUSR1, syscall.SIGUSR2,
		syscall.SIGTERM)
	for s := range c {
		fmt.Printf("Signals  get  %s ........... \n", s)
		break
	}

}

var logger zap.Logger

func logTest(title string) {

	for i := 0; ; i++ {
		logger.Warn(
			"Or use strongly-typed wrappers to add structured context.",
			zap.String("title", title),
			zap.Int("line", i),
		)
		//time.Sleep(time.Second * 1)
		time.Sleep(time.Millisecond * 100)
	}

}

var loggerMutex = new(sync.Mutex)

func main() {

	f, err := os.Open("./logzap.txt")

	if err != nil {
		fmt.Println(err)
		f, err = os.Create("./logzap.txt")
		fmt.Println(err)
	}

	//return

	defer f.Close()

	logger = zap.New(
		//zap.NewJSONEncoder(zap.NoTime()), // drop timestamps in tests

		zap.NewJSONEncoder(zap.TimeFormatter(func(t time.Time) zap.Field {
			return zap.String("time", t.String())
		})),
		zap.Output(f),

		zap.Hook(func(t *zap.Entry) error {
			if t.Time.Second() == 0 {
				fmt.Println("t.Time.Second():", t.Time.Second())
			}
			return nil
		}),
	)

	go logTest("a")
	go logTest("b")
	go logTest("c")
	handleSignal()
}

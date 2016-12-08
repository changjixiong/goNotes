package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
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

var loggerJson zap.Logger
var loggerText zap.Logger

func logJson(title string) {

	for i := 0; ; i++ {

		loggerJson.Warn(
			"Or use strongly-typed wrappers to add structured context.",
			zap.String("title", title),
			zap.Int("line", i),
		)

		time.Sleep(time.Millisecond * 100)
	}

}

func logText(title string) {

	for i := 0; ; i++ {

		loggerText.Log(zap.InfoLevel, "Info1", zap.Int("ID", 1))
		loggerText.Info("Info2")

		loggerText.Log(zap.WarnLevel, "Warn1")
		loggerText.Warn("Warn2", zap.Int("ID", 1))

		loggerText.Log(zap.DebugLevel, "Debug1")
		loggerText.Debug("Debug2")

		loggerText.Log(zap.ErrorLevel, "Error1")
		loggerText.Error("Error2")

		time.Sleep(time.Millisecond * 100)
	}

}

func openLogFile(zapLogPath string) (logFile *os.File) {

	logDir, _ := filepath.Abs(filepath.Dir(zapLogPath))

	if _, err := os.Stat(logDir); err != nil {
		os.Mkdir(logDir, os.ModePerm)
	}

	if _, err := os.Stat(zapLogPath); os.IsNotExist(err) {
		logFile, err = os.Create(zapLogPath)
		if err != nil {
			fmt.Println("Create", zapLogPath, "error:", err)
			return nil
		}

	} else {
		logFile, err = os.OpenFile(zapLogPath, os.O_APPEND|os.O_RDWR, os.ModeAppend)

		if err != nil {
			fmt.Println("Open", zapLogPath, "error:", err)
			return nil
		}
	}

	return logFile
}

func main() {

	f := openLogFile("./log/logzap.txt")
	fText := openLogFile("./log/logzapText.txt")

	defer f.Close()
	defer fText.Close()

	loggerJson = zap.New(

		zap.NewJSONEncoder(zap.TimeFormatter(func(t time.Time) zap.Field {
			return zap.String("time", t.String())
		})),
		zap.Output(f),
		/*
			zap.Hook(func(t *zap.Entry) error {
				if t.Time.Second() == 0 {
					fmt.Println("t.Time.Second():", t.Time.Second())
				}
				return nil
			}),
		*/
	)

	loggerText = zap.New(
		//zap.NewTextEncoder(...)
		zap.NewTextEncoder(zap.TextTimeFormat("2006-01-02 15:04:05")),
		//zap.DebugLevel,
		zap.Output(fText),
	)

	go logJson("a")
	go logText("b")
	handleSignal()
}

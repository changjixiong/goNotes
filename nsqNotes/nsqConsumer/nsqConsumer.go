package main

import (
	"encoding/json"
	"fmt"
	"goNotes/reflectinvoke"

	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	//"github.com/bitly/go-nsq"
	nsq "github.com/nsqio/go-nsq"
)

var reflectinvoker *reflectinvoke.Reflectinvoker

type Foo struct {
}

type Bar struct {
}

func (b *Bar) BarFuncAdd(argOne, argTwo float64) float64 {

	return argOne + argTwo
}

func (f *Foo) FooFuncSwap(argOne, argTwo string) (string, string) {

	return argTwo, argOne
}

func HandleJsonMessage(message *nsq.Message) error {

	resultJson := reflectinvoker.InvokeByJson([]byte(message.Body))
	result := reflectinvoke.Response{}
	err := json.Unmarshal(resultJson, &result)
	if err != nil {
		return err
	}
	info := "HandleJsonMessage get a result\n"
	info += "raw:\n" + string(resultJson) + "\n"
	info += "function: " + result.FuncName + " \n"
	info += fmt.Sprintf("result: %v\n", result.Data)
	info += fmt.Sprintf("error: %d,%s\n\n", result.ErrorCode,
		reflectinvoke.ErrorMsg(result.ErrorCode))

	fmt.Println(info)

	return nil
}

func HandleStringMessage(message *nsq.Message) error {

	fmt.Printf("HandleStringMessage get a message  %v\n\n", string(message.Body))
	return nil
}

func MakeConsumer(topic, channel string, config *nsq.Config,
	handle func(message *nsq.Message) error) {
	consumer, _ := nsq.NewConsumer(topic, channel, config)
	consumer.AddHandler(nsq.HandlerFunc(handle))

	// 待深入了解
	// 連線到 NSQ 叢集，而不是單個 NSQ，這樣更安全與可靠。
	// err := q.ConnectToNSQLookupd("127.0.0.1:4161")

	err := consumer.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Panic("Could not connect")
	}
}

func init() {

	foo := &Foo{}
	bar := &Bar{}

	reflectinvoker = reflectinvoke.NewReflectinvoker()
	reflectinvoker.RegisterMethod(foo)
	reflectinvoker.RegisterMethod(bar)
}

func main() {

	config := nsq.NewConfig()
	config.DefaultRequeueDelay = 0
	config.MaxBackoffDuration = 20 * time.Millisecond
	config.LookupdPollInterval = 1000 * time.Millisecond
	config.RDYRedistributeInterval = 1000 * time.Millisecond
	config.MaxInFlight = 2500

	MakeConsumer("Topic_string", "ch", config, HandleStringMessage)
	MakeConsumer("Topic_json", "ch", config, HandleJsonMessage)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	fmt.Println("exit")

}

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"goNotes/dbnotes/dbhelper"
	"goNotes/dbnotes/model"
	cassandra "goNotes/dbnotes/modelcassandra"
)

var wg sync.WaitGroup

var whichDB = flag.String("whichDB", "mysql", "mysql or cassandra")
var dbIP = flag.String("dbIP", "127.0.0.1", "db ip")
var dbPort = flag.Int("dbPort", 3306, "db port")
var dbName = flag.String("dbName", "dbnote", "db name")
var dbUser = flag.String("dbUser", "root", "db user")
var dbPass = flag.String("dbPass", "123456", "db password")
var keyspace = flag.String("keyspace", "space_for_back", "cassandra Keyspace")

var lastNum10ms = flag.Int("lastNum10ms", -1, "lastNum10ms")
var lastNum100ms = flag.Int("lastNum100ms", -1, "lastNum100ms")
var lastNum2s = flag.Int("lastNum2s", -1, "lastNum2s")
var lastNum4s = flag.Int("lastNum4s", -1, "lastNum4s")

func mysqlFunc() {
	dbhelper.Init(*dbIP, *dbPort, *dbName, *dbUser, *dbPass)
	notices, _ := model.DefaultNotice.QueryByMap(map[string]interface{}{"1": "1"})

	if len(notices) > 0 {
		// 测试代码，更新第0行记录
		notices[0].Content = notices[0].Content + "update"
		notices[0].Update()
	} else {
		emoji := string([]byte{240, 159, 152, 143})
		(&model.Notice{
			SenderID:   123,
			ReceiverID: 234,
			// Content:    "new",
			Content:    emoji,
			Createtime: time.Now(),
			Status:     0}).Insert()
	}

	msgs, _ := model.DefaultMsg.QueryByMap(map[string]interface{}{"content": "def"})
	mails, _ := model.DefaultMail.QueryByMap(map[string]interface{}{"Title": "t1"})

	if len(msgs) > 0 {
		msgs[0].Delete()
	}

	if len(mails) > 0 {
		mails[0].Delete()
	}

	msg := &model.Msg{SenderID: 123,
		ReceiverID: 234,
		Content:    "abc",
		Createtime: time.Now(),
		Status:     0}
	msg.Insert()

	msg = &model.Msg{SenderID: 123,
		ReceiverID: 234,
		Content:    "def",
		Createtime: time.Now(),
		Status:     0}

	msg.Insert()

	mail := &model.Mail{SenderID: 123,
		ReceiverID: 234,
		Title:      "t1",
		Content:    "abc",
		Createtime: time.Now(),
		Status:     0}

	mail.Insert()

	mail = &model.Mail{SenderID: 123,
		ReceiverID: 234,
		Title:      "t2",
		Content:    "abc",
		Createtime: time.Now(),
		Status:     0}

	mail.Insert()

	msgs, _ = model.DefaultMsg.QueryByMap(map[string]interface{}{"content": "def"})
	mails, _ = model.DefaultMail.QueryByMap(map[string]interface{}{"Title": "t1"})

	msgs[0].Content = "update"
	msgs[0].Update()

	mails[0].Content = "update"
	mails[0].Update()

	for _, m := range msgs {
		fmt.Println(m)
	}

	for _, m := range mails {
		fmt.Println(m)
	}

	fmt.Println("OK")

}

func insertCass(tableName string, lastNum int) {
	defer wg.Done()

	f := func(numBegin int) {}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var dur time.Duration

	switch tableName {
	case "num_log_10ms":
		dur = time.Millisecond * 10
		f = func(numBegin int) {
			cassandra.NumLog10msOp.Insert(
				&cassandra.NumLog10ms{
					CreateTime: time.Now(),
					Num:        numBegin,
					ServerID:   r.Intn(4) + 1,
				},
			)
		}
	case "num_log_100ms":
		dur = time.Millisecond * 100
		f = func(numBegin int) {
			cassandra.NumLog100msOp.Insert(
				&cassandra.NumLog100ms{
					CreateTime: time.Now(),
					Num:        numBegin,
					ServerID:   r.Intn(4) + 1,
				},
			)
		}
	case "num_log_2s":
		dur = time.Second * 2
		f = func(numBegin int) {
			cassandra.NumLog2sOp.Insert(
				&cassandra.NumLog2s{
					CreateTime: time.Now(),
					Num:        numBegin,
					ServerID:   r.Intn(4) + 1,
				},
			)
		}
	case "num_log_4s":
		dur = time.Second * 4
		f = func(numBegin int) {
			cassandra.NumLog4sOp.Insert(
				&cassandra.NumLog4s{
					CreateTime: time.Now(),
					Num:        numBegin,
					ServerID:   r.Intn(4) + 1,
				},
			)
		}
	}

	for !cancelled() {
		lastNum += 1
		f(lastNum)
		time.Sleep(dur)
	}

	fmt.Println(tableName, "lastNum:", lastNum)
}

func cassandraFunc() {
	dbhelper.InitCassDB([]string{*dbIP}, *dbPort, *keyspace, "Quorum", *dbUser, *dbPass)

	if *lastNum10ms >= 0 {
		wg.Add(1)
		go insertCass("num_log_10ms", *lastNum10ms)
	}

	if *lastNum100ms >= 0 {
		wg.Add(1)
		go insertCass("num_log_100ms", *lastNum100ms)
	}

	if *lastNum2s >= 0 {
		wg.Add(1)
		go insertCass("num_log_2s", *lastNum2s)
	}

	if *lastNum4s >= 0 {
		wg.Add(1)
		go insertCass("num_log_4s", *lastNum4s)
	}

	handleSignal()

	wg.Wait()
	dbhelper.CloseCassDB()
}

func handleSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	for s := range c {
		fmt.Println("get signal", s, "stopping ......")
		break
	}

	close(done)
}

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

var done = make(chan struct{})

func main() {

	flag.Parse()

	switch *whichDB {
	case "mysql":
		mysqlFunc()
	case "cassandra":
		cassandraFunc()
	default:
		log.Fatal("unknow db")
	}

}

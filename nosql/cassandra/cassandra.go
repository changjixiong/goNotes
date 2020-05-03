package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gocql/gocql"
)

/*
CREATE KEYSPACE IF NOT EXISTS example
WITH REPLICATION = {'class': 'SimpleStrategy','replication_factor':1};

CREATE TABLE example.msg
( id uuid,content varchar,create_time timestamp,PRIMARY KEY (id));
*/

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

func insertMsg(session *gocql.Session) {
	for {

		if err := session.Query(`INSERT INTO msg (id, content, create_time) VALUES (?, ?, ?)`,
			gocql.TimeUUID(),
			"info", time.Now()).Exec(); err != nil {
			fmt.Println("insertMsg err:", err)
		}

		time.Sleep(time.Second * 1)
	}
}
func insertMsgBatch(session *gocql.Session) {

	for {
		batch := session.NewBatch(gocql.LoggedBatch)
		for i := 0; i < 10; i++ {

			batch.Query(`INSERT INTO msg (id, content, create_time) 
				VALUES (?, ?, ?)`,
				gocql.TimeUUID(),
				"info", time.Now())

		}

		if err := session.ExecuteBatch(batch); err != nil {
			fmt.Println("execute batch:", err)
		}

		time.Sleep(time.Second * 3)
	}

}

func main() {

	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()

	if nil != err {
		fmt.Println("CreateSession err:", err)
		return
	}

	defer session.Close()

	go insertMsg(session)
	go insertMsgBatch(session)

	handleSignal()
}

package main

import (
	"fmt"

	"goNotes/dbnotes/model"
)

func main() {

	msg := &model.Msg{SenderID: 123,
		ReceiverID: 234,
		Content:    "abc",
		Status:     0}

	model.DefaultMsg.Insert(msg)

	msg = &model.Msg{SenderID: 123,
		ReceiverID: 234,
		Content:    "def",
		Status:     0}

	model.DefaultMsg.Insert(msg)

	mail := &model.Mail{SenderID: 123,
		ReceiverID: 234,
		Title:      "t1",
		Content:    "abc",
		Status:     0}

	model.DefaultMail.Insert(mail)

	mail = &model.Mail{SenderID: 123,
		ReceiverID: 234,
		Title:      "t2",
		Content:    "abc",
		Status:     0}

	model.DefaultMail.Insert(mail)

	msgs, _ := model.DefaultMsg.QueryByMap(map[string]interface{}{"content": "def"})
	mails, _ := model.DefaultMail.QueryByMap(map[string]interface{}{"Title": "t1"})

	for _, m := range msgs {
		fmt.Println(m)
	}

	for _, m := range mails {
		fmt.Println(m)
	}

	fmt.Println("OK")

}

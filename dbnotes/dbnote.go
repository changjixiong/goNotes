package main

import (
	"fmt"

	"goNotes/dbnotes/dbhelper"
	"goNotes/dbnotes/model"
)

func main() {

	dbhelper.Init("127.0.0.1", 3306, "dbnote", "root", "123456")

	notices, _ := model.DefaultNotice.QueryByMap(map[string]interface{}{"1": "1"})

	if len(notices) > 0 {
		notices[0].Content = notices[0].Content + "update"
		notices[0].Update()
	} else {
		emoji := string([]byte{240, 159, 152, 143})
		(&model.Notice{
			No:         0,
			SenderID:   123,
			ReceiverID: 234,
			// Content:    "new",
			Content: emoji,
			Status:  0}).Insert()
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
		Status:     0}
	msg.Insert()

	msg = &model.Msg{SenderID: 123,
		ReceiverID: 234,
		Content:    "def",
		Status:     0}

	msg.Insert()

	mail := &model.Mail{SenderID: 123,
		ReceiverID: 234,
		Title:      "t1",
		Content:    "abc",
		Status:     0}

	mail.Insert()

	mail = &model.Mail{SenderID: 123,
		ReceiverID: 234,
		Title:      "t2",
		Content:    "abc",
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

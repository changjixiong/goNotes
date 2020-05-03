package main

import (
	"bytes"
	"fmt"
	. "goNotes/dnsnotes/dnsKit"
	"net"
)

func main() {

	// convertTest()

	// queryTest1()

	dnsMsg1 := DNSMessage{
		Header: &DNSHeader{
			ID:                  0x0010,
			QR:                  0,
			OperationCode:       0,
			AuthoritativeAnswer: 0,
			Truncation:          0,
			RecursionDesired:    1,
			RecursionAvailable:  0,
			Zero:                0,
			ResponseCode:        0,
			QuestionCount:       1,
			AnswerRRs:           0,
			AuthorityRRs:        0,
			AdditionalRRs:       0,
		},
		Questions: []*DNSQuestion{
			&DNSQuestion{
				QuestionName:  "www.test1.com",
				QuestionType:  1,
				QuestionClass: 1,
			},
		},
	}

	dnsMsg2 := DNSMessage{
		Header: &DNSHeader{
			ID:                  0x0010,
			QR:                  0,
			OperationCode:       0,
			AuthoritativeAnswer: 0,
			Truncation:          0,
			RecursionDesired:    1,
			RecursionAvailable:  0,
			Zero:                0,
			ResponseCode:        0,
			QuestionCount:       1,
			AnswerRRs:           0,
			AuthorityRRs:        0,
			AdditionalRRs:       0,
		},
		Questions: []*DNSQuestion{
			&DNSQuestion{
				QuestionName:  "www.test2.com",
				QuestionType:  1,
				QuestionClass: 1,
			},
		},
	}

	dnsServer := "127.0.0.1:11153"
	var conn net.Conn
	var err error
	if conn, err = net.Dial("udp", dnsServer); err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	if _, err := conn.Write(dnsMsg1.ToBytes()); err != nil {
		fmt.Println(err.Error())
		return
	}
	buf := make([]byte, 1024)

	if length, err := conn.Read(buf); err == nil {
		result := NewDNSMessage(bytes.NewBuffer(buf[0:length]))
		fmt.Println("query:", result.Questions[0].QuestionName, ", get ip:", result.ResourceRecodes[0].RData)
	} else {
		fmt.Println(err.Error())
	}

	if _, err = conn.Write(dnsMsg2.ToBytes()); err != nil {
		fmt.Println(err.Error())
		return
	}

	buf = make([]byte, 1024)

	if length, err := conn.Read(buf); err == nil {
		result := NewDNSMessage(bytes.NewBuffer(buf[0:length]))
		fmt.Println("query:", result.Questions[0].QuestionName, ", get ip:", result.ResourceRecodes[0].RData)
	} else {
		fmt.Println(err.Error())
	}

}

func queryTest1() {

	dnsmsg := DNSMessage{
		Header: &DNSHeader{
			// ID: 0x0001,
			ID:                  0x0010,
			QR:                  0,
			OperationCode:       0,
			AuthoritativeAnswer: 0,
			Truncation:          0,
			RecursionDesired:    1,
			RecursionAvailable:  0,
			Zero:                0,
			ResponseCode:        0,
			QuestionCount:       1,
			AnswerRRs:           0,
			AuthorityRRs:        0,
			AdditionalRRs:       0,
		},
		Questions: []*DNSQuestion{
			&DNSQuestion{
				QuestionName:  "www.baidu.com",
				QuestionType:  1,
				QuestionClass: 1,
			},
		},
	}

	var conn net.Conn
	var err error
	// dnsServer := "114.114.114.114:53"
	dnsServer := "127.0.0.1:11153"
	if conn, err = net.Dial("udp", dnsServer); err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	if _, err := conn.Write(dnsmsg.ToBytes()); err != nil {
		fmt.Println(err.Error())
		return
	}
	buf := make([]byte, 1024)
	length, err := conn.Read(buf)

	if err == nil {
		fmt.Println(buf[0:length])
	} else {
		fmt.Println(err.Error())
	}
}

func convertTest() {

	header := &DNSHeader{
		ID: 1,

		QR:                  2,
		OperationCode:       3,
		AuthoritativeAnswer: 4,
		Truncation:          5,
		RecursionDesired:    6,
		RecursionAvailable:  7,
		Zero:                8,
		ResponseCode:        9,

		QuestionCount: 10,
		AnswerRRs:     11,
		AuthorityRRs:  12,
		AdditionalRRs: 13,
	}
	fmt.Println(header.ToBytes())
	fmt.Println(NewDNSHeader(bytes.NewBuffer(header.ToBytes())).ToBytes())

	question := &DNSQuestion{
		QuestionName:  "www.baidu.com",
		QuestionType:  1,
		QuestionClass: 1,
	}
	fmt.Println(question.ToBytes())
	fmt.Println(NewDNSQuestion(bytes.NewBuffer(question.ToBytes())).ToBytes())

}

package dnsKit

import (
	"bytes"
	"encoding/binary"
	"net"
	"strings"
)

/*
DNS报文格式，不论是请求报文，还是DNS服务器返回的应答报文，都使用统一的格式。
2个字节(16bit)，标识字段，客户端会解析服务器返回的DNS应答报文，获取ID值与请求报文设置的ID值做比较
如果相同，则认为是同一个DNS会话。
QR	标示该消息是请求消息（该位为0）还是应答消息（该位为1）
QR这一段位 Flag 16bit长度
header
  0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                      ID                       |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|QR|  opcode   |AA|TC|RD|RA|   Z    |   RCODE   |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    QDCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    ANCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    NSCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    ARCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

question
  0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                     ...                       |
|                    QNAME                      |
|                     ...                       |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    QTYPE                      |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    QCLASS                     |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/

// DNSHeader 头
type DNSHeader struct {
	ID uint16

	//golang 没有bit类型，只能这样处理
	QR                  uint16 //1bit
	OperationCode       uint16 //4bit
	AuthoritativeAnswer uint16 //1bit
	Truncation          uint16 //1bit
	RecursionDesired    uint16 //1bit
	RecursionAvailable  uint16 //1bit
	Zero                uint16 //3bit
	ResponseCode        uint16 //4bit

	QuestionCount uint16
	AnswerRRs     uint16
	AuthorityRRs  uint16
	AdditionalRRs uint16
}

func (dh *DNSHeader) ToBytes() []byte {
	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, dh.ID)

	bits := dh.QR<<15 + dh.OperationCode<<11 + dh.AuthoritativeAnswer<<10 + dh.Truncation<<9
	bits += dh.RecursionDesired<<8 + dh.RecursionAvailable<<7 + dh.ResponseCode
	binary.Write(&buffer, binary.BigEndian, bits)

	binary.Write(&buffer, binary.BigEndian, dh.QuestionCount)
	binary.Write(&buffer, binary.BigEndian, dh.AnswerRRs)
	binary.Write(&buffer, binary.BigEndian, dh.AuthorityRRs)
	binary.Write(&buffer, binary.BigEndian, dh.AdditionalRRs)

	return buffer.Bytes()
}

func NewDNSHeader(buffer *bytes.Buffer) *DNSHeader {
	// 读12byte
	id := binary.BigEndian.Uint16(buffer.Next(2))
	flag := binary.BigEndian.Uint16(buffer.Next(2))
	return &DNSHeader{
		ID: id,

		QR:                  flag >> 15,
		OperationCode:       (flag >> 11) % (1 << 4),
		AuthoritativeAnswer: (flag >> 10) % (1 << 1),
		Truncation:          (flag >> 9) % (1 << 1),
		RecursionDesired:    (flag >> 8) % (1 << 1),
		RecursionAvailable:  (flag >> 7) % (1 << 1),
		Zero:                (flag >> 4) % (1 << 3),
		ResponseCode:        flag % (1 << 4),

		QuestionCount: binary.BigEndian.Uint16(buffer.Next(2)),
		AnswerRRs:     binary.BigEndian.Uint16(buffer.Next(2)),
		AuthorityRRs:  binary.BigEndian.Uint16(buffer.Next(2)),
		AdditionalRRs: binary.BigEndian.Uint16(buffer.Next(2)),
	}
}

//DNSQuestion 请求
type DNSQuestion struct {
	QuestionName  string
	QuestionType  uint16
	QuestionClass uint16
}

func (dq *DNSQuestion) ToBytes() []byte {
	var buffer bytes.Buffer

	segments := strings.Split(dq.QuestionName, ".")
	for _, seg := range segments {
		binary.Write(&buffer, binary.BigEndian, byte(len(seg)))
		binary.Write(&buffer, binary.BigEndian, []byte(seg))
	}
	binary.Write(&buffer, binary.BigEndian, byte(0x00))

	binary.Write(&buffer, binary.BigEndian, dq.QuestionType)
	binary.Write(&buffer, binary.BigEndian, dq.QuestionClass)

	return buffer.Bytes()
}

func NewDNSQuestion(buffer *bytes.Buffer) *DNSQuestion {
	//8bit标记每一级域名的长度
	// buf := bytes.NewBuffer(data)
	length := uint8(0)
	binary.Read(buffer, binary.BigEndian, &length)
	segments := []string{}
	for length > 0 {
		seg := make([]byte, length)
		binary.Read(buffer, binary.BigEndian, &seg)
		segments = append(segments, string(seg))
		binary.Read(buffer, binary.BigEndian, &length)
		// fmt.Println(length)
	}

	question := &DNSQuestion{
		QuestionName: strings.Join(segments, "."),
	}

	question.QuestionType = binary.BigEndian.Uint16(buffer.Next(2))
	question.QuestionClass = binary.BigEndian.Uint16(buffer.Next(2))

	return question
}

//DNSResourceRecode 回答字段，授权字段，附加字段
type DNSResourceRecode struct {
	Name     string
	NamePos  uint16
	RRType   uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    string
}

func (drr *DNSResourceRecode) ToBytes() []byte {
	var buffer bytes.Buffer
	if drr.NamePos > 0 {
		binary.Write(&buffer, binary.BigEndian, (0x01<<15)|(0x01<<14)|drr.NamePos)
	} else {
		segments := strings.Split(drr.Name, ".")
		for _, seg := range segments {
			binary.Write(&buffer, binary.BigEndian, byte(len(seg)))
			binary.Write(&buffer, binary.BigEndian, []byte(seg))
		}
		binary.Write(&buffer, binary.BigEndian, byte(0x00))
	}

	binary.Write(&buffer, binary.BigEndian, drr.RRType)
	binary.Write(&buffer, binary.BigEndian, drr.Class)
	binary.Write(&buffer, binary.BigEndian, drr.TTL)
	binary.Write(&buffer, binary.BigEndian, drr.RDLength)
	if drr.Class == 1 {
		binary.Write(&buffer, binary.BigEndian, []byte(net.ParseIP(drr.RData).To4()))
	} else if drr.Class == 5 {

	}

	return buffer.Bytes()
}

func NewDNSResourceRecode(buffer *bytes.Buffer) *DNSResourceRecode {
	tag := buffer.Next(1)[0] >> 6
	buffer.UnreadByte()
	drr := &DNSResourceRecode{}
	if tag == 3 {
		//最高两位11，右移后是3
		drr.NamePos = (binary.BigEndian.Uint16(buffer.Next(2)) << 2) >> 2
	} else {

	}
	drr.RRType = binary.BigEndian.Uint16(buffer.Next(2))
	drr.Class = binary.BigEndian.Uint16(buffer.Next(2))
	drr.TTL = binary.BigEndian.Uint32(buffer.Next(4))
	drr.RDLength = binary.BigEndian.Uint16(buffer.Next(2))

	if drr.RRType == 1 && drr.RDLength == 4 {
		drr.RData = net.IPv4(buffer.Next(1)[0], buffer.Next(1)[0], buffer.Next(1)[0], buffer.Next(1)[0]).String()
	} else {
		//FIXME:
		//域名处理
	}

	return drr
}

// DNSMessage 消息体
type DNSMessage struct {
	Header          *DNSHeader
	Questions       []*DNSQuestion
	ResourceRecodes []*DNSResourceRecode
}

func (dm *DNSMessage) ToBytes() []byte {
	result := []byte{}
	result = append(result, dm.Header.ToBytes()...)

	for i := uint16(0); i < dm.Header.QuestionCount; i++ {
		result = append(result, dm.Questions[i].ToBytes()...)
	}

	for _, rr := range dm.ResourceRecodes {
		result = append(result, rr.ToBytes()...)
	}

	return result
}

func NewDNSMessage(buffer *bytes.Buffer) *DNSMessage {
	//用bytes.Buffer类型来逐个字节读取后处理的优点就是不需要自己计算读取偏移值
	//这个对于Question和Answer这种第一段长度不固定的内容处理非常方便
	dnsMsg := &DNSMessage{
		Header: NewDNSHeader(buffer),
	}

	//FIXME: deal more then 1 QuestionCount
	dnsMsg.Questions = append(dnsMsg.Questions, NewDNSQuestion(buffer))

	//FIXME: deal more then 1 ResourceRecode
	if buffer.Len() > 0 {
		dnsMsg.ResourceRecodes = append(dnsMsg.ResourceRecodes, NewDNSResourceRecode(buffer))
	}

	return dnsMsg
}

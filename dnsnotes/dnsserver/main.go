package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	. "goNotes/dnsnotes/dnsKit"
	"net"
	"sync"
)

const RECV_BUF_LEN = 1024

var addr2IP Address2IP

type Address2IP struct {
	lastIP uint32 //167772160 -> 10.0.0.0
	sync.RWMutex
	address2ip map[string]uint32
}

func (a *Address2IP) getIP(address string) string {
	a.RLock()
	ip := uint32(0)
	ok := false
	if ip, ok = a.address2ip[address]; ok {
		a.RUnlock()
	} else {
		a.RUnlock()

		a.Lock()
		a.lastIP++
		a.address2ip[address] = a.lastIP
		ip = a.lastIP
		a.Unlock()

	}

	ipByte := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(ipByte, ip)
	return net.IPv4(ipByte[0], ipByte[1], ipByte[2], ipByte[3]).String()
}

func main() {
	addr2IP = Address2IP{
		lastIP:     167772160,
		address2ip: map[string]uint32{},
	}
	udpAddr, err := net.ResolveUDPAddr("udp", ":11153")
	if err != nil {
		panic("error listening:" + err.Error())
	}
	fmt.Println("Starting the server")

	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()

	for {

		buf := make([]byte, RECV_BUF_LEN)
		n, raddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			panic("Error accept:" + err.Error())
		}
		fmt.Println("Accepted the Connection :", conn.RemoteAddr())
		go echoServer(conn, raddr, buf[0:n])
	}
}

func echoServer(conn *net.UDPConn, raddr *net.UDPAddr, data []byte) {

	dnsMsg := NewDNSMessage(bytes.NewBuffer(data))

	dnsMsg.Header.QR = 1
	dnsMsg.Header.RecursionAvailable = 1
	dnsMsg.Header.AnswerRRs = 1

	dnsMsg.ResourceRecodes = append(
		dnsMsg.ResourceRecodes,
		&DNSResourceRecode{
			Name:     "",
			NamePos:  12,
			RRType:   1,
			Class:    1,
			TTL:      384,
			RDLength: 4,
			RData:    addr2IP.getIP(dnsMsg.Questions[0].QuestionName),
		})

	fmt.Println("query:", dnsMsg.Questions[0].QuestionName, ", ip:", dnsMsg.ResourceRecodes[0].RData)
	conn.WriteToUDP(dnsMsg.ToBytes(), raddr)

}

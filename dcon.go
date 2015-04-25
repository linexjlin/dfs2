package main

import (
	"bufio"
	"net"
	"sync"
	"time"
)

const (
	PING = 0xFF
	PONG = 0xFF
)

//这个类实现了，tcp连接保持功能。 这里我们通过C/S ping pong 的方式实现了tcp保持，并检测连接。
type DCon struct {
	con     net.Conn
	lc      sync.Mutex
	cmsgOut chan HostMsg
}

func (dc *DCon) sendPing() error {
	//	fmt.Println("befor Send")
	_, e := dc.con.Write([]byte{PING})
	//	fmt.Println("after send")
	return e
}

func (dc *DCon) receivePong() error {
	//	fmt.Println("befor receive")
	b := make([]byte, 1)
	_, e := dc.con.Read(b)
	//	fmt.Println("after receive")
	//	fmt.Println("receive", n, b)
	return e
}

func (dc *DCon) receive(b []byte) (n int, err error) {
	dc.lc.Lock()
	dc.con.SetDeadline(time.Now().Add(time.Second * 10))
	nn, ee := dc.con.Read(b)
	dc.con.SetDeadline(time.Time{})
	dc.lc.Unlock()
	return nn, ee
}

func (dc *DCon) ReceiveMsgSend() (err error) {
	dc.lc.Lock()
	dc.con.SetDeadline(time.Now().Add(time.Second * 10))
	reader := bufio.NewReader(dc.con)
	buf, _, e := reader.ReadLine()
	dc.con.SetDeadline(time.Time{})
	dc.lc.Unlock()
	dc.con.RemoteAddr().String()
	if e == nil {
		dc.cmsgOut <- HostMsg{dc.con.RemoteAddr().String(), string(buf)}
	}
	return e
}

func (dc *DCon) send(b []byte) (n int, err error) {
	dc.lc.Lock()
	dc.con.SetDeadline(time.Now().Add(time.Second * 10))
	nn, ee := dc.con.Write(b)
	dc.con.SetDeadline(time.Time{})
	dc.lc.Unlock()
	return nn, ee
}

//This function to keep connection between C/S
func (dc *DCon) KeepConn() net.Conn {
	for {
		dc.lc.Lock()
		dc.con.SetDeadline(time.Now().Add(time.Second * 20))
		if e := dc.sendPing(); e != nil {
			//			fmt.Println(e)
			dc.con.Close()
			return dc.con
		}
		if e := dc.receivePong(); e != nil {
			//			fmt.Println(e)
			dc.con.Close()
			return dc.con
		}
		dc.con.SetDeadline(time.Time{})
		dc.lc.Unlock()
		time.Sleep(time.Second * 30)
	}
	return dc.con
}

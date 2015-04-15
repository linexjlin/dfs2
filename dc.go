package main

import (
	"fmt"
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
	con net.Conn
	lc  sync.Mutex
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
		time.Sleep(time.Second * 3)
	}
	return dc.con
}

//用来维护所有客户端连接。 在个类里实现了，保持tcp连接，断开连接自动清理功能。
type DCons struct {
	DCmap map[DCon]int
	lc    sync.Mutex
	count int
}

func NewDCons() *DCons {
	tmap := make(map[DCon]int)
	var l sync.Mutex
	ndcs := &DCons{tmap, l, 0}
	ndcs.count = 0
	return ndcs
}

func (dcs *DCons) AddToKeep(c net.Conn) {
	var tlc sync.Mutex
	dc := DCon{c, tlc}
	dcs.DCmap[dc] = dcs.count
	dcs.count++
	go func() {
		dcs.Del(dc.KeepConn())
	}()
	return
}

func (dcs *DCons) Del(c net.Conn) {
	var tlc sync.Mutex
	c.Close()
	tdc := DCon{c, tlc}
	delete(dcs.DCmap, tdc)
	dcs.count--
	return
}

func (dcs *DCons) DelAll() {
	dcs.lc.Lock()
	for dc, _ := range dcs.DCmap {
		delete(dcs.DCmap, dc)
		dcs.count--
	}
	dcs.count = 0
	dcs.lc.Unlock()
	return
}

func (dcs *DCons) List() {
	fmt.Println("There are ", dcs.count, "socket(s)")
	for dc, _ := range dcs.DCmap {
		fmt.Println(dc.con.RemoteAddr())
	}
}

//在所有连接里广播消息
func (dcs *DCons) Broadcast(msg string) {
	for dc, _ := range dcs.DCmap {

		fmt.Println("the romote one is:", dc.con.RemoteAddr())
		dc.send([]byte(msg))
	}
}

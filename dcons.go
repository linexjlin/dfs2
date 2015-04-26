package main

import (
	"fmt"
	"net"
	"sync"
)

//用来维护所有客户端连接。 在个类里实现了，保持tcp连接，断开连接自动清理功能。
type DCons struct {
	DCmap   map[DCon]int
	lc      sync.Mutex
	count   int
	gmsgOut chan HostMsg
}

func NewDCons() *DCons {
	//	tmap := make(map[DCon]int)
	//	var l sync.Mutex
	ndcs := &DCons{}
	ndcs.count = 0
	return ndcs
}

func (dcs *DCons) Init() {
	dcs.count = 0
	dcs.DCmap = make(map[DCon]int)
}

func (dcs *DCons) AddToKeep(c net.Conn) {
	var tlc sync.Mutex
	ch := dcs.gmsgOut
	dc := DCon{c, tlc, ch}
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
	tdc := DCon{c, tlc, dcs.gmsgOut}
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

package main

import (
	"net"
	"time"
)

var nameMap map[string]chan net.Conn

//server routine
//func serverSide() {
//	server := NewDServer()
//	server.SetListenAddr(":8612")
//	server.Listen()
//}

/*
//the client wait receive message from server routine
func clientSide() {
	var ccons = NewDCons()
	var com ConCom
	serverHosts := []string{"127.0.0.1:8612", "localhost:8612"}
	for _, h := range serverHosts {
		c, _ := net.Dial("tcp", h)
		e := com.SendMyType(c, []byte("client"))
		if e != nil {
			c.Close()
		}
		ccons.AddToKeep(c)
	}
	var com ConCom
	//Todo: wait to redial
	for {
		time.Sleep(time.Second * 3)
		fmt.Println(len(ccons.DCmap))
	}
}*/

//http routine
func httpSide() {
	for {
		time.Sleep(time.Second * 3)
	}
}

func main() {
	var nrc = make(chan NameReader)
	dh := DHttp{}
	dh.dc.Nrc = nrc
	dh.Init()
	go dh.HttpListen()

	s := DServer{}
	s.Init()
	go s.ServerListen()
	s.Nrc = nrc
	go s.ReceiveName()
	go s.FileListen()

	c := DClient{}
	c.Init()
	c.Dial()
	c.ReceiveName()
	for {
		time.Sleep(time.Second * 2)
	}

}

// testdcon project main.go
package main

import (
	"fmt"
	"net"
	"time"
)

var nameMap map[string]chan net.Conn

//server routine
func serverSide() {
	var scons = NewDCons()
	var com ConCom
	scons.count = 0
	l, _ := net.Listen("tcp", ":8612")
	for {
		s, _ := l.Accept()
		defer s.Close()
		//time.Sleep(time.Second * 10)
		//			getConType()
		ctype, e := com.GetConType(s)
		if e != nil {
			s.Close()
		}
		fmt.Println("receive：", ctype)
		switch ctype {
		case "client":
			scons.AddToKeep(s)
			scons.List()
		default:
			s.Close()
		}
		//			scons.Broadcast("Here comes new firend!")

		//	serveConn(s)
	}
}

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
	//Todo: wait to redial
	for {
		time.Sleep(time.Second * 3)
		fmt.Println(len(ccons.DCmap))
	}
}

//http routine
func httpSide() {
	for {
		time.Sleep(time.Second * 3)
	}
}

func main() {
	go serverSide()
	go clientSide()
	httpSide()
}

// testdcon project main.go
package main

import (
	"fmt"
	"net"
	"time"
)

var scons = NewDCons()
var ccons = NewDCons()
var nameMap map[string]chan net.Conn

//server routine
func serverSide() {
	scons.count = 0
	l, _ := net.Listen("tcp", ":8612")
	for {
		s, _ := l.Accept()
		//time.Sleep(time.Second * 10)
		//			getConType()
		scons.AddToKeep(s)
		scons.List()
		//			scons.Broadcast("Here comes new firend!")
		defer s.Close()
		//	serveConn(s)
	}
}

//client routine
func clientSide() {
	serverHosts := []string{"127.0.0.1:8612", "localhost:8612"}
	for _, h := range serverHosts {
		c, _ := net.Dial("tcp", h)
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

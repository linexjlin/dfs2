package main

import (
	"net"
	"time"
)
import "io"
import lg "fmt"
import "os"

type DServer struct {
	dcons      DCons
	laddr      string
	fileLsnadr string
	fileReader map[string]chan io.Reader
	clientCnt  int
	Nrc        chan NameReader
}

//func (ds *DServer) init() {
//	ds.clientCnt = 0
//	ds.dcons.count = 0
//	return
//}

func (ds *DServer) AddName(nc NameReader) {
	if _, ok := ds.fileReader[nc.Name]; ok {
		return
	}
	ds.fileReader[nc.Name] = nc.Rchan
	lg.Println("already have", len(ds.fileReader), "request")
	go ds.autoClean(nc.Name)
	return
}

func (ds *DServer) autoClean(fn string) {
	time.Sleep(time.Second * 5)

	delete(ds.fileReader, fn)
	lg.Println("auto clean ", fn, "complete")
}

//func (ds *DServer) SetFileListenAddr(addr string) error {
//	ds.fileLsnadr = addr
//	return nil
//}

//func (ds *DServer) SetListenAddr(addr string) error {
//	ds.laddr = addr
//	return nil
//}

func (ds *DServer) ServerListen() {
	lg.Println("Client server listen on:", ds.laddr)
	l, e := net.Listen("tcp", ds.laddr)
	if e != nil {
		lg.Println(e)
		os.Exit(1)
	}
	for {
		s, _ := l.Accept()
		defer s.Close()
		ds.clientCnt++
		lg.Println("New Client comming:", s.RemoteAddr())
		ds.dcons.AddToKeep(s)
	}
}

func (ds *DServer) FileListen() {
	lg.Println("File receiver listen on:", ds.fileLsnadr)
	l, e := net.Listen("tcp", ds.fileLsnadr)
	if e != nil {
		lg.Println(e)
		os.Exit(1)
	}
	for {
		s, _ := l.Accept()
		defer s.Close()
		go ds.DealFileCon(s)
	}
}

func (ds *DServer) DealFileCon(s net.Conn) {
	var con ConCom
	fn, _ := con.GetString(s)
	lg.Println("The reader of", fn, "comes")
	ds.fileReader[fn] <- s
	delete(ds.fileReader, fn)
}

func (ds *DServer) GetReader(fileName string) io.Reader {
	select {
	case reader := <-ds.fileReader[fileName]:
		return reader
	}
	return nil
}

func (ds *DServer) ReceiveName() {
	for {
		select {
		case nc := <-ds.Nrc:
			go func() {
				lg.Println("server get new request, file name:", nc.Name)
				ds.AddName(nc)
				lg.Println("Broadcast start, wait to get reader from", ds.dcons.count, "client(s)")
				ds.dcons.Broadcast(nc.Name) //then wait client to connect
				//				rd := ds.GetReader(nc.Name)
				//				nc.Rchan <- rd
				//				delete(ds.fileReader, nc.Name)
			}()
		}
	}
}

func (ds *DServer) Init() {
	ds.fileLsnadr = ":8702"
	ds.laddr = ":8612"
	ds.dcons.Init()
	ds.fileReader = make(map[string]chan io.Reader)
}

//func NewDServer() *DServer {
//	d := new(DServer)
//	//	d.init()
//	return d
//}

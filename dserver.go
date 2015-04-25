package main

import "net"
import "io"

type DServer struct {
	dcons      DCons
	laddr      string
	fladdr     string
	fileReader map[string]chan io.Reader
	clientCnt  int
}

func (ds *DServer) init() {
	ds.clientCnt = 0
	ds.dcons.count = 0
	return
}

func (ds *DServer) SetFileListenAddr(addr string) error {
	ds.fladdr = addr
	return nil
}

func (ds *DServer) SetListenAddr(addr string) error {
	ds.laddr = addr
	return nil
}

func (ds *DServer) Listen() {
	l, _ := net.Listen("tcp", ds.laddr)
	for {
		s, _ := l.Accept()
		defer s.Close()
		ds.dcons.AddToKeep(s)
	}
}

func (ds *DServer) FileListen() {
	l, _ := net.Listen("tcp", ds.fladdr)
	for {
		s, _ := l.Accept()
		defer s.Close()
	}
}

func (ds *DServer) GetReader(fileName string) io.Reader {
	select {
	case reader := <-ds.fileReader[fileName]:
		return reader
	}
	return nil
}

func (ds *DServer) ReceiveNamePostReader(msgc <-chan string, rc chan<- io.Reader) {
	for {
		select {
		case fileName := <-msgc:
			ds.dcons.Broadcast(fileName)
			rc <- ds.GetReader(fileName)
			delete(ds.fileReader, fileName)
		}
	}
}

func NewDServer() *DServer {
	d := new(DServer)
	d.init()
	return d
}

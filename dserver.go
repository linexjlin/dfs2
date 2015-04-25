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

func (ds *DServer) AddName(fn string) {
	if _, ok := ds.fileReader[fn]; ok {
		return
	}
	var rd chan io.Reader
	ds.fileReader[fn] = rd
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
		ds.clientCnt++
		ds.dcons.AddToKeep(s)
	}
}

func (ds *DServer) DealFileCon(s net.Conn) {
	var con ConCom
	fn, _ := con.GetString(s)
	ds.fileReader[fn] <- s
}

func (ds *DServer) FileListen() {
	l, _ := net.Listen("tcp", ds.fladdr)
	for {
		s, _ := l.Accept()
		defer s.Close()
		go ds.DealFileCon(s)
	}
}

func (ds *DServer) GetReader(fileName string) io.Reader {
	select {
	case reader := <-ds.fileReader[fileName]:
		return reader
	}
	return nil
}

func (ds *DServer) ReceiveNamePostReader(nameChan chan NameReader) {
	for {
		select {
		case nr := <-nameChan:
			go func() {
				ds.AddName(nr.name)
				ds.dcons.Broadcast(nr.name)
				rd := ds.GetReader(nr.name)
				nameChan <- NameReader{nr.name, rd}
				delete(ds.fileReader, nr.name)
			}()

		}
	}
}

func NewDServer() *DServer {
	d := new(DServer)
	d.init()
	return d
}

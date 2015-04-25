package main

import (
	"io"
	"time"
)
import "net"
import "os"

type FReader struct {
	fileName string
	reader   chan io.Reader
}

type DClient struct {
	dcons      DCons
	serverList []string
	paths      []string
	serverCnt  int
	fReaderc   chan FReader
}

func (dc *DClient) Dial() {
	for _, h := range dc.serverList {
		c, _ := net.Dial("tcp", h)
		defer c.Close()
		dc.dcons.AddToKeep(c)
	}
}

func (dc *DClient) ReceiveFileName() {
	for {
		select {
		case hm := <-dc.dcons.gmsgOut:
			go dc.ClientDoing(hm)
		}
	}
}

//Find reader from local first, if no found ask outside(server) for reader
func (dc *DClient) GetReader(fileName string) (r io.Reader, finfo os.FileInfo) {
	if r, finfo = dc.FindLocalReader(fileName); r != nil {
		return r, finfo
	} else {
		rc := make(chan io.Reader, 1)
		dc.fReaderc <- FReader{fileName, rc} //ask outside for reader
		reader := dc.GetOutsideReader(rc)
		if reader != nil {
			return reader, nil
		}
		return nil
	}
}

//Get reader from outside
func (dc *DClient) GetOutsideReader(rc chan io.Reader) io.Reader {
	select {
	case reader := <-rc:
		return reader
	case <-time.After(time.Second * 3):
		return nil
	}
}

func (dc *DClient) GetWriter(host string) io.Writer {
	con, e := net.Dial("tcp", host)
	if e != nil {
		return nil
	}
	return con
}

//
func (dc *DClient) ClientDoing(hm HostMsg) {
	fileName := hm.msg
	host := hm.host

	//try to get reader
	reader finfo:= dc.GetReader(fileName)
	if reader == nil {
		return nil
	}

	//get upper writer
	writer := dc.GetWrier(host)
	if writer == nil {
		return nil
	}
	if finfo!=nil {
		ds.WriteInfo(finfo,writer)
		io.Copy(writer,reader)
	}
}


func (dc *DClient) WriteInfo(finfo os.FileInfo,w io.Writer) error {
	return nil
}

func (dc *DClient) FindLocalReader(fileName string) (r io.Reader, finfo os.FileInfo) {
	var file *os.File
	for _, path := range dc.paths {
		info, e := os.Stat(path + fileName)
		if e == nil {
			file, _ = os.Open(path + fileName)
			finfo = info
			break
		}
	}
	return file, finfo
}

func (dc *DClient) SendFileInfo(fileName string) io.Reader {
	return nil
}

func (dc *DClient) GetTransCon(fileName string) io.Reader {
	return nil
}

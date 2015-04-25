package main

import (
	"io"
	"time"
)
import "net"
import "os"

type DClient struct {
	dcons      DCons
	serverList []string
	paths      []string
	serverCnt  int
	nrc        chan NameReader //name reader chan
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
		var rc io.Reader
		dc.nrc <- NameReader{fileName, rc} //ask outside for reader
		reader := dc.GetOutsideReader()
		if reader != nil {
			return reader, nil
		}
		return nil, nil
	}
}

//Get reader from outside
func (dc *DClient) GetOutsideReader() io.Reader {
	select {
	case nr := <-dc.nrc:
		return nr.reader
	case <-time.After(time.Second * 3):
		return nil
	}
}

func (dc *DClient) GetWriter(hn HostMsg) io.Writer {
	con, e := net.Dial("tcp", hn.host)
	if e != nil {
		return nil
	}
	var c ConCom
	_, e2 := c.PutString(hn.msg, con)
	if e2 != nil {
		return nil
	}
	return con
}

//
func (dc *DClient) ClientDoing(hm HostMsg) {
	fileName := hm.msg

	//try to get reader
	reader, finfo := dc.GetReader(fileName)
	if reader == nil {
		return
	}

	//get upper writer
	writer := dc.GetWriter(hm)
	if writer == nil {
		return
	}
	if finfo != nil {
		dc.WriteInfo(finfo, writer)
		io.Copy(writer, reader)
	}
}

func (dc *DClient) WriteInfo(finfo os.FileInfo, w io.Writer) error {
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

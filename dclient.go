package main

import (
	"bytes"
	"errors"
	lg "fmt"
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
	Nrc        chan NameReader //name reader chan
	MaxWait    time.Duration
}

func (dc *DClient) Dial() {
	for _, h := range dc.serverList {
		lg.Println("dialing", h)
		c, e := net.Dial("tcp", h)
		if e != nil {
			lg.Println(e)
			return
		}
		lg.Println("Success dialed", h)
		lg.Println("local one is:", c.LocalAddr())
		defer c.Close()
		dc.dcons.AddToKeep(c)
	}
}
func (dc *DClient) Init() {
	dc.paths = []string{"."}
	dc.serverList = []string{"127.0.0.1:8612"}
	dc.serverCnt = 0
	dc.dcons.Init()
}

func (dc *DClient) ReceiveName() {
	for {
		select {
		case hm := <-dc.dcons.gmsgOut:
			lg.Println("get message from", hm.host, ". The message is ", hm.msg)
			go dc.ClientDoing(hm)
		}
	}
}

//Find reader from local first, if no found ask outside(server) for reader
func (dc *DClient) GetReader(fileName string) (r io.Reader, err error) {

	if rr, ee := dc.FindLocalReader(fileName); ee == nil {
		return rr, ee
	} else {

		var rc chan io.Reader
		lg.Println("local can not find, try outside")

		//		dc.Nrc <- NameReader{fileName, rc} //ask outside for reader
		//send message to server side
		go func() {
			//			gnrc <- NameReader{fileName, rc}
			dc.Nrc <- NameReader{fileName, rc} //ask outside for reader
		}()
		reader, e := dc.GetOutsideReader(rc)
		return reader, e
	}
}

func (dc *DClient) FindLocalReader(fileName string) (r io.Reader, err error) {
	var file *os.File
	rw := new(bytes.Buffer)
	for _, path := range dc.paths {
		info, e := os.Stat(path + fileName)
		if e == nil {
			file, _ = os.Open(path + fileName)
			defer file.Close()
			finfo := Finfo{}
			finfo.IsDir = info.IsDir()
			finfo.ModTime = info.ModTime()
			finfo.Name = info.Name()
			finfo.Size = info.Size()
			finfo.GobEncodeToWriter(rw)
			io.Copy(rw, file)
			return rw, e
			break
		}
	}
	return nil, errors.New("No local reader foud!")
}

//Get reader from outside
func (dc *DClient) GetOutsideReader(rc chan io.Reader) (io.Reader, error) {
	lg.Println("Wait", dc.MaxWait, "to get outside reader!")
	select {
	case rd := <-rc:
		return rd, nil
	case <-time.After(dc.MaxWait):
		return nil, errors.New("Not found outside reader!")
	}
}

func (dc *DClient) GetWriter(hn HostMsg) (io.Writer, error) {
	con, e := net.Dial("tcp", hn.host)
	if e != nil {
		return nil, e
	}
	var c ConCom
	_, e2 := c.PutString(hn.msg, con)
	if e2 != nil {
		return nil, e2
	}
	return con, e2
}

//
func (dc *DClient) ClientDoing(hm HostMsg) {
	fileName := hm.msg

	//try to get reader
	reader, e := dc.GetReader(fileName)
	if e != nil {
		return
	}

	//get upper writer
	writer, ee := dc.GetWriter(hm)
	if ee != nil {
		return
	}

	//copy reader to writer
	io.Copy(writer, reader)

}

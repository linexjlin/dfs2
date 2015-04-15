package main

import (
	"errors"
	"io"
	"net/http"
)
import "net"
import "time"

const (
	NOFOUND = 0x00
)

type ConnInfo struct {
	size    int64
	dtype   string
	found   bool
	transed bool
	crtTime time.Time
}

type HTC struct {
	nc        chan net.Conn
	rw        http.ResponseWriter
	tcpWriter io.Reader
	tcpReader io.Reader
}

//文件接收服务端类型
type DFileHTTP struct {
	htcMap map[string]*HTC
	//	fileMap map[string]*ConnInfo //这里只能放指针，
}

func (df *DFileHTTP) AddName(fileName string) error {
	//添加名字的时候注意不要重名，重名要等待。
	return nil
}

//等待获取tcpreader在这里最多等待3秒，超过3秒，
func (df *DFileHTTP) GetReader(fileName string) (r io.Reader, e error) {
	cnt := 0
	timeout := make(chan int, 1)

	go func() {
		select {
		case t := <-time.After(time.Second * 3):
			timeout <- 1
		}
	}()

	//在这里等待获取，文件流的reader
	for {
		select {
		case c := <-df.htcMap[fileName].nc:
			if isFound(c) {
				return c, nil
			}
			cnt++
			c.Close()
		case <-timeout:
			return nil, errors.New("Time out and no Reader found")
		}
	}
	return nil, errors.New("interal error!")
}

//Todo 这个函数不用写在这里
func (df *DFileHTTP) ServeTCPFile(rw http.ResponseWriter, r *http.Request) {
	//	lg.Println("Current URL is:", r.URL.Path)

	if r.URL.Path == "/" {
		http.NotFound(rw, r)
		return
	}

	//broadcast 用文件指针做
	fileName := r.URL.Path
	df.htcMap[fileName].rw = rw
	var e error
	df.htcMap[fileName].tcpReader, e = df.GetReader(fileName)
	if e == nil {
		io.Copy(df.htcMap[fileName].rw, df.htcMap[fileName].tcpReader)
	} else {
		http.NotFound(rw, r)
	}
}

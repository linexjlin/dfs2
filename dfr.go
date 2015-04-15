package main

import "net/http"
import "net"
import "time"

type ConnInfo struct {
	size        int64
	dtype       string
	transCon    net.Conn
	fileInfoCon net.Conn
	found       bool
	transed     bool
	crtTime     time.Time
}

type DFileHTTP struct {
	chanMap map[string]chan net.Conn
	fileMap map[string]*ConnInfo //这里只能放指针，
}

func (df *DFileHTTP) AddName(fileName string) error {
	//添加名字的时候注意不要重名，重名要等待。
	return nil
}

func (df *DFileHTTP) waitResponse(fileName string) bool {
	cnt := 0
	timeout := make(chan int, 1)

	go func() {
		select {
		case t := <-time.After(time.Second * 3):
			timeout <- 1
		}
	}()

	for {
		select {
		case c := <-df.chanMap[fileName]:
			if isFound(c) {
				df.fileMap[fileName].fileInfoCon = c
				return true
			}
			cnt++
			c.Close()
		case <-timeout:
			return false
		}
	}
	return false
}

func (df *DFileHTTP) ServeTCPFile(rw http.ResponseWriter, r *http.Request) {
	//	lg.Println("Current URL is:", r.URL.Path)

	if r.URL.Path == "/" {
		http.NotFound(rw, r)
		return
	}

	//	fileName := getFileName(r.URL.Path)
	//	registFilename()
	//	broadcast() //todo 函数指针
	found := df.waitResponse(fileName)
	if found {
		//		getFileInfo() //todo
		//		setHttpHeader() //todo
		tcpFileToHTTP(rw, fileReceiver)
	} else {
		http.NotFound(rw, r)
	}
}

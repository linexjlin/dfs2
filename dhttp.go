package main

import (
	lg "fmt"
	"io"
	"net/http"
	"time"
)

type DHttp struct {
	dc   DClient
	port string
}

func (dh *DHttp) HttpHandle(rw http.ResponseWriter, r *http.Request) {
	dh.HttpDoing(rw, r)
}

func (dh *DHttp) HttpListen() {
	h := http.HandlerFunc(dh.HttpHandle)
	lg.Println("HTTP listen on:", dh.port)
	http.ListenAndServe(":"+dh.port, h)
}

func (dh *DHttp) SetFileInfo(rw http.ResponseWriter) {
	return
}

func (dh *DHttp) HttpDoing(w http.ResponseWriter, r *http.Request) error {
	fileName := r.URL.Path
	lg.Println("Get fileName", fileName)

	//try to get reader
	lg.Println("Wait to get reader")
	reader, e := dh.dc.GetReader(fileName)

	if e != nil {
		http.NotFound(w, r) //Can't not find any file
		return nil
	}
	var info Finfo
	info.GobDecodeFromReader(reader)
	lg.Println("Parsed size is:", info.Size)

	dh.SetFileInfo(w)
	io.Copy(w, reader)
	return nil
}

func (dh *DHttp) Init() {
	dh.port = "8701"
	//	dh.dc.paths = []string{"."}
	dh.dc.MaxWait = time.Second * 1
}

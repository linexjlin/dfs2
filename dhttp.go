package main

import (
	"io"
	"net/http"
)

type DHttp struct {
	dc   DClient
	port string
}

func (dh *DHttp) HttpHandle(rw http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path
	dh.HttpDoing(fileName, rw)
}

func (dh *DHttp) HttpListen() {

}

func (dh *DHttp) GetFileInfo(reader io.Reader) {

}

func (dh *DHttp) SetFileInfo(rw http.ResponseWriter) {

}

func (dh *DHttp) HttpDoing(fn string, w http.ResponseWriter) error {
	fileName := fn

	//try to get reader
	reader, finfo := dh.dc.GetReader(fileName)
	if reader == nil {
		return nil
	}
	dh.GetFileInfo(reader)
	dh.SetFileInfo(w)

	io.Copy(w, reader)
	return nil
}

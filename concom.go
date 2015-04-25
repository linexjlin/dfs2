package main

import (
	"bufio"
	"errors"
	"net"
)

//这一类型用来在socket间建立简单的通信
const RSP = 0x01

type ConCom struct {
}

//判断是否收到来连接的响应
func (com ConCom) receiveResponse(con net.Conn) bool {
	b := make([]byte, 1)
	_, e := con.Read(b)
	if e != nil {
		return false
	}
	return true
}

func (com ConCom) sendResponse(con net.Conn) bool {
	_, e := con.Write([]byte{RSP})
	if e != nil {
		return false
	}
	return true
}

//向服务器发送这个连接的类型。
func (com ConCom) SendMyType(con net.Conn, myType []byte) error {
	con.Write([]byte(string(myType) + "\n"))
	if com.receiveResponse(con) {
		return nil
	}
	return errors.New("Can't send file type, socket may closed")
}

func (com ConCom) GetConType(con net.Conn) (ctype string, err error) {
	reader := bufio.NewReader(con)
	rst, _, e := reader.ReadLine()
	com.sendResponse(con)
	return string(rst), e
}

func (com ConCom) GetString(con net.Conn) (str string, err error) {
	reader := bufio.NewReader(con)
	rst, _, e := reader.ReadLine()
	return string(rst), e
}

func (com ConCom) PutString(str string, con net.Conn) (n int, err error) {
	nn, e := con.Write([]byte(str))
	return nn, e
}

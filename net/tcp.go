package net

import (
	"net"
	"time"
)

// DialTCP dials the tcp [addr] with [timeout] and returns the [conn]
func DialTCP(addr string, timeout int) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*time.Duration(timeout))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ListenTCP(addr string, client chan net.Conn) {
	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		Fail2Resolve(addr)
	}
	ln, err := net.ListenTCP("tcp", address)
	if err != nil {
		Fail2Listen(addr)
	}
	defer ln.Close()
	Listening(address.String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			Fail2Accept(err)
		}
		Accepted(conn.RemoteAddr().String())
		client <- conn
	}
}

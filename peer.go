package main

import (
	"net"
)

type Peer struct {
	con net.Conn
	msgChan chan []byte
}

func newPeer(connection net.Conn, msgChan chan []byte) *Peer{
	return &Peer{
		con : connection,
		msgChan : msgChan,
	}
}

func (p *Peer) readRequest() error{
	buff := make([]byte,1024)

	for {
		n, err := p.con.Read(buff)
		if err != nil{
			return err
		}
		msgBuff := make([]byte, n)
		copy(msgBuff,buff)
		p.msgChan <- msgBuff
	}
	
}
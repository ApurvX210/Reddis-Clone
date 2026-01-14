package peers

import (
	"net"
)

type Message struct {
	Data []byte
	Peer *Peer
}

type Peer struct {
	con net.Conn
	msgChan chan Message
}

func (p *Peer) Send(msg []byte) (int,error){
	return p.con.Write(msg)
}

func NewPeer(connection net.Conn, msgChan chan Message) *Peer{
	return &Peer{
		con : connection,
		msgChan : msgChan,
	}
}

func (p *Peer) ReadRequest() error{
	buff := make([]byte,1024)

	for {
		n, err := p.con.Read(buff)
		if err != nil{
			return err
		}
		msgBuff := make([]byte, n)
		copy(msgBuff,buff)
		p.msgChan <- Message{msgBuff,p}
	}
	
}
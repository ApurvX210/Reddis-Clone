package main

import (
	"fmt"
	"log/slog"
	"net"
)

const defaultListenerAddress = ":5001"

type Config struct {
	ListenerAddress string
}

type Server struct {
	Config
	peers       map[*Peer]bool
	ln          net.Listener
	addPeerChan chan *Peer
	quitChan    chan struct{}
	msgChan     chan []byte
}

func newServer(cfg Config) *Server {
	if len(cfg.ListenerAddress) == 0 {
		cfg.ListenerAddress = defaultListenerAddress
	}
	return &Server{
		Config:      cfg,
		peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
		quitChan:    make(chan struct{}),
		msgChan:     make(chan []byte),
	}
}

func (s *Server) start() error {
	slog.Info("Server Running","Listning on Port ",s.ListenerAddress)
	ln, err := net.Listen("tcp", s.ListenerAddress)
	if err != nil {
		return err
	}
	s.ln = ln
	go s.chanLoop()
	return s.accecptRequest()
}

func (s *Server) accecptRequest() error {
	for {
		con, err := s.ln.Accept()
		if err != nil {
			slog.Error(fmt.Sprintf("Error occured while accepting request-> %s", err))
			continue
		}
		go s.handleConnection(con)
	}
}

func (s *Server) handleRawMsg(rawMsg []byte) error{

}

func (s *Server) chanLoop() {
	for {
		select {
			case rawMsg := <-s.msgChan:
				if err := s.handleRawMsg(rawMsg); err != nil{
					slog.Error("Error Occured while Handling Raw message ","Error",err)
				}
				fmt.Println(string(rawMsg))
			case peer := <-s.addPeerChan:
				s.peers[peer] = true
			case <-s.quitChan:
				return
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	peer := newPeer(conn,s.msgChan)
	s.addPeerChan <- peer
	slog.Info("New Peer connected","remoteAddress",conn.RemoteAddr)
	if err := peer.readRequest(); err != nil {
		slog.Error("Error Occured while reading Peer request ","err",err,"RemoteAddress",conn.RemoteAddr)
	}
}

func main() {
	server := newServer(Config{ListenerAddress: "127.0.0.1:5000"})
	fmt.Println(server.start())
}

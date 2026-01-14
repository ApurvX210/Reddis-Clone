package main

import (
	"REDDIS/parsing"
	"REDDIS/peers"
	"REDDIS/storage"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/tidwall/resp"
)

const defaultListenerAddress = ":5000"

type Config struct {
	ListenerAddress string
}

type Server struct {
	Config
	peers       map[*peers.Peer]bool
	ln          net.Listener
	addPeerChan chan *peers.Peer
	delPeerChan chan *peers.Peer
	quitChan    chan struct{}
	msgChan     chan peers.Message
	db          *storage.DB
}

func newServer(cfg Config) *Server {
	if len(cfg.ListenerAddress) == 0 {
		cfg.ListenerAddress = defaultListenerAddress
	}
	return &Server{
		Config:      cfg,
		peers:       make(map[*peers.Peer]bool),
		addPeerChan: make(chan *peers.Peer),
		delPeerChan: make(chan *peers.Peer),
		quitChan:    make(chan struct{}),
		msgChan:     make(chan peers.Message, 1024),
		db:          storage.NewDb(),
	}
}

func (s *Server) start() error {
	slog.Info("Server Running", "Listning on Port ", s.ListenerAddress)
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

func (s *Server) handleMsg(message peers.Message) error {
	cmd, err := parsing.ParseCommand(string(message.Data))
	if err != nil {
		return err
	}
	switch cmd := cmd.(type) {
	case parsing.SetCommand:
		err = s.db.Set(cmd.Key, cmd.Value)
		var msg []byte
		if err != nil {
			msg = []byte(fmt.Sprintf("Error occured while executing set command key: %s value: %s", cmd.Key, cmd.Value))
		} else {
			msg = resp.SimpleStringValue("OK").Bytes()
		}
		_, err = message.Peer.Send(msg)
		if err != nil {
			return err
		}
	case parsing.GetCommand:
		val, response := s.db.Get(cmd.Key)
		var msg []byte
		if !response {
			msg = []byte(fmt.Sprintf("Key %s not found", cmd.Key))
		} else {
			msg = resp.AnyValue(val).Bytes()
		}
		_, err = message.Peer.Send(msg)
		if err != nil {
			return err
		}
	case parsing.DelCommand:
		s.db.Get(cmd.Key)
		var msg []byte
		msg = []byte(fmt.Sprintf("Key %s deleted successfully", cmd.Key))
		_, err = message.Peer.Send(msg)
		if err != nil {
			return err
		}
	case parsing.HelloCommad:
		response := s.db.Hello()
		_, err = message.Peer.Send([]byte(response))
		if err != nil {
			return err
		}
	case parsing.ClientInfoCommand:
		response := resp.SimpleStringValue("OK").Bytes()
		_, err = message.Peer.Send(response)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) chanLoop() {
	for {
		select {
		case message := <-s.msgChan:
			if err := s.handleMsg(message); err != nil {
				slog.Error("Error Occured while Handling Raw message ", "Error", err)
			}
		case peer := <-s.addPeerChan:
			s.peers[peer] = true
		case peer := <-s.delPeerChan:
			delete(s.peers, peer)
		case <-s.quitChan:
			return
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	peer := peers.NewPeer(conn, s.msgChan)
	s.addPeerChan <- peer
	slog.Info("New Peer connected", "remoteAddress", conn.RemoteAddr)
	if err := peer.ReadRequest(); err != nil {
		if err == io.EOF {
			slog.Info("Peer Connection closed ", "remoteAddress", conn.RemoteAddr)
			s.delPeerChan <- peer
		} else {
			slog.Error("Error Occured while reading Peer request ", "err", err, "RemoteAddress", conn.RemoteAddr)
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"github.com/tidwall/resp"
)
const defaultListenerAddress = ":5000"

type Config struct {
	ListenerAddress string
}

type Message struct {
	data []byte
	peer *Peer
}

type Server struct {
	Config
	peers       map[*Peer]bool
	ln          net.Listener
	addPeerChan chan *Peer
	delPeerChan chan *Peer
	quitChan    chan struct{}
	msgChan     chan Message
	db          *DB
}

func newServer(cfg Config) *Server {
	if len(cfg.ListenerAddress) == 0 {
		cfg.ListenerAddress = defaultListenerAddress
	}
	return &Server{
		Config:      cfg,
		peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
		delPeerChan: make(chan *Peer),
		quitChan:    make(chan struct{}),
		msgChan:     make(chan Message, 1024),
		db:          NewDb(),
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

func (s *Server) handleMsg(message Message) error {
	cmd, err := parseCommand(string(message.data))
	if err != nil {
		return err
	}
	switch cmd := cmd.(type) {
	case SetCommand:
		err = s.db.Set(cmd.key, cmd.value)
		var msg []byte
		if err != nil {
			msg = []byte(fmt.Sprintf("Error occured while executing set command key: %s value: %s", cmd.key, cmd.value))
		} else {
			msg = resp.SimpleStringValue("OK").Bytes()
		}
		_, err = message.peer.Send(msg)
		if err != nil {
			return err
		}
	case GetCommand:
		val, response := s.db.Get(cmd.key)
		var msg []byte
		if !response {
			msg = []byte(fmt.Sprintf("Key %s not found", cmd.key))
		} else {
			msg = resp.AnyValue(val).Bytes()
		}
		_, err = message.peer.Send(msg)
		if err != nil {
			return err
		}
	case HelloCommad:
		fmt.Println("Apurv Hello")
		response := s.db.hello()
		_, err = message.peer.Send([]byte(response))
		if err != nil {
			return err
		}
	case ClientInfoCommand:
		fmt.Println("Apurv")
		response :=  resp.SimpleStringValue("OK").Bytes()
		_, err = message.peer.Send(response)
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
			delete(s.peers,peer)
		case <-s.quitChan:
			return
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	peer := newPeer(conn, s.msgChan)
	s.addPeerChan <- peer
	slog.Info("New Peer connected", "remoteAddress", conn.RemoteAddr)
	if err := peer.readRequest(); err != nil {
		if err == io.EOF{
			slog.Info("Peer Connection closed ", "remoteAddress", conn.RemoteAddr)
			s.delPeerChan <- peer
		}else{
			slog.Error("Error Occured while reading Peer request ", "err", err, "RemoteAddress", conn.RemoteAddr)
		}
	}
}

func main() {
	listenAddress := flag.String("listenAddress",defaultListenerAddress,"Listen Adress of reddis server")
	flag.Parse()
	fmt.Println(*listenAddress)
	server := newServer(Config{ListenerAddress: *listenAddress})
	log.Fatal(server.start())
}

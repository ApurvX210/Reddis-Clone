package main

import (
	"REDDIS/client"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"
)

const defaultListenerAddress = ":5001"

type Config struct {
	ListenerAddress string
}

type Message struct{
	data []byte
	peer *Peer
}

type Server struct {
	Config
	peers       map[*Peer]bool
	ln          net.Listener
	addPeerChan chan *Peer
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
		quitChan:    make(chan struct{}),
		msgChan:     make(chan Message,1024),
		db: 		 NewDb(),
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

func (s *Server) handleMsg(message Message) error{
	cmd,err := parseCommand(string(message.data))
	if err != nil{
		return err
	}
	switch cmd := cmd.(type){
		case SetCommand:
			err =  s.db.Set(cmd.key,cmd.value)
			var msg string
			if err != nil{
				msg = fmt.Sprintf("Error occured while executing set command key: %s value: %s",cmd.key,cmd.value)
			}else{
				msg = fmt.Sprintf("Successfully executed set command key: %s value: %s",cmd.key,cmd.value)
			}
			_,err = message.peer.Send([]byte(msg))
			if err != nil{
				return err
			}
		case GetCommand:
			val,response:= s.db.Get(cmd.key)
			var msg string
			if !response{
				msg = fmt.Sprintf("Key %s not found",cmd.key)
			}else{
				msg = string(val)
			}
			_,err = message.peer.Send([]byte(msg))
			if err != nil{
				return err
			}
	}
	return nil
}

func (s *Server) chanLoop() {
	for {
		select {
			case message := <-s.msgChan:
				if err := s.handleMsg(message); err != nil{
					slog.Error("Error Occured while Handling Raw message ","Error",err)
				}
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
	go func ()  {
		fmt.Println(server.start())
	}()
	time.Sleep(time.Second)
	cl,err := client.New("localhost:5000")

	if err != nil{
		log.Fatal(err)
	}

	if response,err := cl.Set(context.Background(),"admin","Apurv"); err!= nil{
		log.Fatal(err)
	}else{
		fmt.Println(response)
	}
	// if err := cl.Set(context.Background(),"heelo","yash"); err!= nil{
	// 	log.Fatal(err)
	// }
	// if err := cl.Set(context.Background(),"ap","sac"); err!= nil{
	// 	log.Fatal(err)
	// }
	// if err := cl.Set(context.Background(),"cascac","cac"); err!= nil{
	// 	log.Fatal(err)
	// }
	// if err := cl.Set(context.Background(),"qwwqw","cac"); err!= nil{
	// 	log.Fatal(err)
	// }
	
	time.Sleep(time.Second*10)
	fmt.Println(server.db.data)

	if response,err := cl.Get(context.Background(),"admin"); err!= nil{
		log.Fatal(err)
	}else{
		fmt.Println(response)
	}

	select{}
	
}

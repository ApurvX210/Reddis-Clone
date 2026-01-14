package main

import (
	"flag"
	"fmt"
	"log"
)


func main() {
	listenAddress := flag.String("listenAddress", defaultListenerAddress, "Listen Adress of reddis server")
	flag.Parse()
	fmt.Println(*listenAddress)
	server := newServer(Config{ListenerAddress: *listenAddress})
	log.Fatal(server.start())
}

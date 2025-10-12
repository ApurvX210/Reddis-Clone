package main

import (
	"fmt"
	"testing"
)

func TestProtocol(t *testing.T) {
	msg := "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	cmd,err := parseCommand(msg)
	if err != nil{
		t.Fatal("Errrorr")
	}
	fmt.Println(cmd)
}

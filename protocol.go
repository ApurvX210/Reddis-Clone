package main

import "fmt"

type Command struct {
}

func parseCommand(msg string) {
	fmt.Println(msg[0])
	
}

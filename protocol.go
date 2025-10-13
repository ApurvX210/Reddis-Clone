package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"github.com/tidwall/resp"
)

const (
	CommandSET = "SET"
	CommandGET = "GET"
	CommandDEL = "DEL"
)

type Command interface {
}

type SetCommand struct {
	key,value []byte
}

type GetCommand struct {
	key []byte
}

type DelCommand struct {
	key []byte
}

func parseCommand(rawMsg string) (Command,error){
	rd := resp.NewReader(bytes.NewBufferString(rawMsg))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error occured while parsing the request","err",err)
		}
		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String(){
					case CommandSET:
						if len(v.Array()) != 3{
							return nil,fmt.Errorf("invalid set comand provided: Invalid no of variables")
						}
						cmd := SetCommand {
							key: v.Array()[1].Bytes(),
							value: v.Array()[2].Bytes(),
						}
						return cmd,nil
					case CommandGET:
						if len(v.Array()) != 2{
							return nil,fmt.Errorf("invalid get comand provided: Invalid no of variables")
						}
						cmd := GetCommand {
							key: v.Array()[1].Bytes(),
						}
						return cmd,nil
					}
					
			}
		}
	}
	return nil,fmt.Errorf("invalid or unknown command revieved: %s",rawMsg)
}

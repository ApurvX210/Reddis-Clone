package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"

	"github.com/tidwall/resp"
)

const (
	CommandSET = "SET"
	CommandGET = "GET"
	CommandDEL = "DEL"
)

type Command interface {
	processCommand()
}

type SetCommand struct {
	key string
	value any
}

type GetCommand struct {
	key string
}

type DelCommand struct {
	key string
}

func (cmd *SetCommand) processCommand(){
	slog.Info("Set Command")
}

func (cmd *GetCommand) processCommand(){
	slog.Info("Get Command")
}

func (cmd *DelCommand) processCommand(){
	slog.Info("Delete Command")
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
						cmd := &SetCommand {
							key: v.Array()[1].String(),
							value: v.Array()[2],
						}
						return cmd,nil
					}
			}
		}
	}
	return nil,fmt.Errorf("invalid or unknown command revieved: %s",rawMsg)
}

package main

import (
	"bytes"
	"fmt"
	"github.com/tidwall/resp"
	"io"
	"log"
	"reflect"
)

const (
	CommandSET   = "set"
	CommandGET   = "get"
	CommandDEL   = "del"
	CommandHELLO = "hello"
	CommandClientInfo = "client"
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

type HelloCommad struct {
	version int
}

type ClientInfoCommand struct {
	libName string
}

func parseCommand(rawMsg string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(rawMsg))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error occured while parsing the request", "err", err)
		}
		fmt.Println(v.Array()[0])
		if v.Type() == resp.Array {
			value := v.Array()[0]
			switch value.String() {
			case CommandSET:
				if len(v.Array()) != 3 {
					return nil, fmt.Errorf("invalid set comand provided: Invalid no of variables")
				}
				cmd := SetCommand{
					key:   v.Array()[1].Bytes(),
					value: v.Array()[2].Bytes(),
				}
				return cmd, nil
			case CommandGET:
				if len(v.Array()) != 2 {
					return nil, fmt.Errorf("invalid get comand provided: Invalid no of variables")
				}
				cmd := GetCommand{
					key: v.Array()[1].Bytes(),
				}
				return cmd, nil
			case CommandHELLO:
				cmd := HelloCommad{
					version: v.Array()[1].Integer(),
				}
				return cmd, nil
			case CommandClientInfo:
				cmd := ClientInfoCommand{
					libName: v.Array()[1].String(),
				}
				return cmd, nil
			}
		}
	}
	return nil, fmt.Errorf("invalid or unknown command revieved: %s", rawMsg)
}

func initialHandShake(m map[string]any) string {
	buf := bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))
	for key, val := range m {
		switch reflect.TypeOf(val).Kind() {
		case reflect.Int:
			buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n:%d\r\n", len(key), key, val))
		default:
			buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n+%s\r\n", len(key), key, val))
		}

	}
	return buf.String()
}

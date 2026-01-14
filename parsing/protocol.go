package parsing

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"

	"github.com/tidwall/resp"
)

const (
	CommandSET        = "set"
	CommandGET        = "get"
	CommandDEL        = "del"
	CommandHELLO      = "hello"
	CommandClientInfo = "client"
)

type Command interface {
}

type SetCommand struct {
	Key, Value []byte
}

type GetCommand struct {
	Key []byte
}

type DelCommand struct {
	Key []byte
}

type HelloCommad struct {
	Version int
}

type ClientInfoCommand struct {
	LibName string
}

func ParseCommand(rawMsg string) (Command, error) {
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
					Key:   v.Array()[1].Bytes(),
					Value: v.Array()[2].Bytes(),
				}
				return cmd, nil
			case CommandGET:
				if len(v.Array()) != 2 {
					return nil, fmt.Errorf("invalid get comand provided: Invalid no of variables")
				}
				cmd := GetCommand{
					Key: v.Array()[1].Bytes(),
				}
				return cmd, nil
			case CommandHELLO:
				cmd := HelloCommad{
					Version: v.Array()[1].Integer(),
				}
				return cmd, nil
			case CommandClientInfo:
				cmd := ClientInfoCommand{
					LibName: v.Array()[1].String(),
				}
				return cmd, nil
			}
		}
	}
	return nil, fmt.Errorf("invalid or unknown command revieved: %s", rawMsg)
}

func InitialHandShake(m map[string]any) string {
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

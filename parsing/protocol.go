package parsing

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"

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
			return nil, fmt.Errorf("error parsing request: %w", err)
		}
		fmt.Println(v.Array()[0])
		if v.Type() == resp.Array {
			value := v.Array()[0]
			// Redis commands are case-insensitive
			cmdStr := strings.ToLower(value.String())
			switch cmdStr {
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
			case CommandDEL:
				if len(v.Array()) != 2 {
					return nil, fmt.Errorf("invalid del comand provided: Invalid no of variables")
				}
				cmd := DelCommand{
					Key: v.Array()[1].Bytes(),
				}
				return cmd, nil
			case CommandHELLO:
				// HELLO can have 0 or more arguments (version is optional)
				version := 2 // Default to RESP2
				if len(v.Array()) >= 2 {
					version = v.Array()[1].Integer()
				}
				cmd := HelloCommad{
					Version: version,
				}
				return cmd, nil
			case CommandClientInfo:
				// CLIENT command can have subcommands like SETNAME, INFO, etc.
				// Handle gracefully - accept any CLIENT command
				libName := ""
				if len(v.Array()) >= 3 {
					// CLIENT SETNAME <name> format
					libName = v.Array()[2].String()
				} else if len(v.Array()) >= 2 {
					// CLIENT <subcommand> or CLIENT <name> format
					libName = v.Array()[1].String()
				}
				cmd := ClientInfoCommand{
					LibName: libName,
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

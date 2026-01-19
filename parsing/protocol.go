package parsing

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
	"github.com/tidwall/resp"
)

const (
	CommandSET        = "set"
	CommandGET        = "get"
	CommandDEL        = "del"
	CommandHELLO      = "hello"
	CommandClientInfo = "client"
	SubCommandSetInfo = "setinfo"
	SubCommandExpiry  = "exp"
	DefaultExpiration = 30
)

type Command interface {
}

type SetCommand struct {
	Key, Value 	[]byte
	Exp			time.Time
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
	info map[string]string
}

func ParseCommand(rawMsg string) (Command, error) {
	fmt.Println(rawMsg)
	rd := resp.NewReader(bytes.NewBufferString(rawMsg))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error parsing request: %w", err)
		}
		// fmt.Println(v.Array()[0])
		if v.Type() == resp.Array {
			value := v.Array()[0]
			// Redis commands are case-insensitive
			cmdStr := strings.ToLower(value.String())
			switch cmdStr {
			case CommandSET:
				if len(v.Array()) < 3 {
					return nil, fmt.Errorf("invalid set comand provided: Invalid no of variables")
				}
				var subCmd string;
				var cmd Command;
				if len(v.Array()) == 5 {
					subCmd = strings.ToLower(v.Array()[3].String())
				} 
				switch subCmd{
				case SubCommandExpiry:
					cmd = SetCommand{
						Key:   v.Array()[1].Bytes(),
						Value: v.Array()[2].Bytes(),
						Exp:   time.Now().Add(time.Duration(v.Array()[4].Float()*float64(time.Second))),
					}
				default:
					cmd = SetCommand{
						Key:   v.Array()[1].Bytes(),
						Value: v.Array()[2].Bytes(),
						Exp:   time.Now().Add(DefaultExpiration*time.Second),
					}
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
				// CLIENT command can have subcommands like SETINFO, SETNAME, INFO, etc.
				// Handle gracefully - accept any CLIENT command
				info := make(map[string]string)

				// Check if we have at least a subcommand
				if len(v.Array()) < 2 {
					cmd := ClientInfoCommand{
						info: info,
					}
					return cmd, nil
				}

				// Subcommands are case-insensitive
				subcmd := strings.ToLower(v.Array()[1].String())

				// Handle CLIENT SETINFO <key> <value> format
				if subcmd == SubCommandSetInfo {
					if len(v.Array()) >= 4 {
						key := v.Array()[2].String()
						value := v.Array()[3].String()
						info[key] = value
					}
				}
				// For other CLIENT subcommands (SETNAME, INFO, etc.), just accept them

				cmd := ClientInfoCommand{
					info: info,
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

package client

import (
	"bytes"
	"context"
	// "fmt"
	"log/slog"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	address string
	conn    net.Conn
}

func New(address string) (*Client,error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil,err
	} else {
		slog.Info("Connection established", "Remote Address", conn.RemoteAddr())
	}
	return &Client{
		address: address,
		conn: conn,
	},nil
}

func (c *Client) Set(ctx context.Context, key string, val any) (string,error) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	// fmt.Println("Set")
	wr.WriteArray([]resp.Value{resp.StringValue("SET"), resp.StringValue(key), resp.AnyValue(val)})
	_, error := c.conn.Write(buf.Bytes())
	if error != nil{
		return "",error
	}
	responseBuff := make([]byte,1024)
	n,err := c.conn.Read(responseBuff)
	return string(responseBuff[:n]),err
}

func (c *Client) Get(ctx context.Context, key string) (string,error) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	// fmt.Println("Get")
	wr.WriteArray([]resp.Value{resp.StringValue("GET"), resp.StringValue(key)})
	_, error := c.conn.Write(buf.Bytes())
	if error != nil{
		return "",error
	}

	responseBuff := make([]byte,1024)
	n,err := c.conn.Read(responseBuff)
	return string(responseBuff[:n]),err
}

func (c *Client) Close() error{
	return c.conn.Close()
}

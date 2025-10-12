package client

import (
	"bytes"
	"context"
	"log/slog"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	address string
}

func New(address string) *Client {
	return &Client{
		address: address,
	}
}

func (c *Client) Set(ctx context.Context, key string, val any) error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}else{
		slog.Info("Connection established","Remote Address",conn.RemoteAddr())
	}

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{resp.StringValue("SET"), resp.StringValue("leader"), resp.StringValue("Charlie")})
	_, error := conn.Write(buf.Bytes())
	return error
}

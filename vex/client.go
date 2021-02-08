package vex

import (
	"bufio"
	"errors"
	"io"
	"net"
)

type Client struct {
	conn net.Conn

	reader io.Reader
}

func NewClient(network string, address string) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}, nil
}

func (c *Client) Do(command byte, args [][]byte) (body []byte, err error) {
	// 包装请求，然后发送给服务端
	_, err = writeRequestTo(c.conn, command, args)
	if err != nil {
		return nil, err
	}

	// 读取服务端的响应
	reply, body, err := readResponseFrom(c.reader)
	if err != nil {
		return body, err
	}

	if reply == ErrorReply {
		return body, errors.New(string(body))
	}
	return body, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

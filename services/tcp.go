package services

import (
	"cache/caches"
	"cache/vex"
	"encoding/binary"
	"encoding/json"
	"errors"
)

const (
	getCommand    = byte(1)
	setCommand    = byte(2)
	deleteCommand = byte(3)
	statusCommand = byte(4)
)

var (
	commandNeedsMoreArgumentsErr = errors.New("command needs more arguments")

	notFoundErr = errors.New("not found")
)

type TCPServer struct {
	cache  *caches.Cache
	server *vex.Server
}

func NewTcpServer(cache *caches.Cache) *TCPServer {
	return &TCPServer{
		cache:  cache,
		server: vex.NewServer(),
	}
}

func (ts *TCPServer) Run(address string) error {
	ts.server.RegisterHandler(getCommand, ts.getHandler)
	ts.server.RegisterHandler(setCommand, ts.setHandler)
	ts.server.RegisterHandler(deleteCommand, ts.deleteHandler)
	ts.server.RegisterHandler(statusCommand, ts.statusHandler)
	return ts.server.ListenAndServer("tcp", address)
}

func (ts *TCPServer) getHandler(args [][]byte) (body []byte, err error) {

	if len(args) < 1 {
		return nil, commandNeedsMoreArgumentsErr
	}

	value, ok := ts.cache.Get(string(args[0]))
	if !ok {
		return value, notFoundErr
	}
	return value, nil
}

func (ts *TCPServer) setHandler(args [][]byte) (body []byte, err error) {
	if len(args) < 3 {
		return nil, commandNeedsMoreArgumentsErr
	}

	ttl := int64(binary.BigEndian.Uint64(args[0]))
	err = ts.cache.SetWithTTL(string(args[1]), args[2], ttl)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ts *TCPServer) deleteHandler(args [][]byte) (body []byte, err error) {
	if len(args) < 1 {
		return nil, commandNeedsMoreArgumentsErr
	}

	err = ts.cache.Delete(string(args[0]))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ts *TCPServer) statusHandler(args [][]byte) (body []byte, err error) {
	return json.Marshal(ts.cache.Status())
}

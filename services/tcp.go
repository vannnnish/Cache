package services

import (
	"cache/caches"
	"cache/helpers"
	"cache/vex"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	getCommand    = byte(1)
	setCommand    = byte(2)
	deleteCommand = byte(3)
	statusCommand = byte(4)
	nodesCommand  = byte(5)
)

var (
	commandNeedsMoreArgumentsErr = errors.New("command needs more arguments")

	notFoundErr = errors.New("not found")
)

type TCPServer struct {
	*node
	cache   *caches.Cache
	server  *vex.Server
	options *Options
}

func NewTcpServer(cache *caches.Cache, options *Options) (*TCPServer, error) {
	n, err := newNode(options)
	if err != nil {
		return nil, err
	}
	return &TCPServer{
		node:    n,
		cache:   cache,
		server:  vex.NewServer(),
		options: options,
	}, nil
}

func (ts *TCPServer) Run() error {
	ts.server.RegisterHandler(getCommand, ts.getHandler)
	ts.server.RegisterHandler(setCommand, ts.setHandler)
	ts.server.RegisterHandler(deleteCommand, ts.deleteHandler)
	ts.server.RegisterHandler(statusCommand, ts.statusHandler)
	ts.server.RegisterHandler(nodesCommand, ts.nodesHandler)
	return ts.server.ListenAndServer("tcp", helpers.JoinAddressAndPort(ts.options.Address, ts.options.Port))
}

func (ts *TCPServer) getHandler(args [][]byte) (body []byte, err error) {

	if len(args) < 1 {
		return nil, commandNeedsMoreArgumentsErr
	}

	// 使用一致性哈希选择出这个 key 所属的物理节点
	key := string(args[0])
	node, err := ts.selectNode(key)
	if err != nil {
		return nil, err
	}

	// 判断这个 key 所属的节点
	if !ts.isCurrentNode(node) {
		return nil, fmt.Errorf("redirect to node %s", node)
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

	// 使用一致性哈希选择出这个 key 所属的物理节点
	key := string(args[1])
	node, err := ts.selectNode(key)
	if err != nil {
		return nil, err
	}

	// 判断这个 key 所属的节点
	if !ts.isCurrentNode(node) {
		return nil, fmt.Errorf("redirect to node %s", node)
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

	// 使用一致性哈希选择出这个 key 所属的物理节点
	key := string(args[0])
	node, err := ts.selectNode(key)
	if err != nil {
		return nil, err
	}

	// 判断这个 key 所属的节点
	if !ts.isCurrentNode(node) {
		return nil, fmt.Errorf("redirect to node %s", node)
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

func (ts *TCPServer) nodesHandler(args [][]byte) (body []byte, err error) {
	return json.Marshal(ts.nodes())
}

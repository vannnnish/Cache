package vex

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"sync"
)

var (
	commandHandlerNotFoundErr = errors.New("failed to find a handler of command")
)

type Server struct {
	// 监听器
	listener net.Listener

	// 命令处理器
	handlers map[byte]func(args [][]byte) (body []byte, err error)
}

func NewServer() *Server {
	return &Server{
		handlers: map[byte]func(args [][]byte) (body []byte, err error){},
	}
}

func (s *Server) RegisterHandler(command byte, handler func(args [][]byte) (body []byte, err error)) {
	s.handlers[command] = handler
}

func (s *Server) ListenAndServer(network string, address string) (err error) {
	s.listener, err = net.Listen(network, address)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.handleConn(conn)
		}()
	}
	wg.Wait()
	return nil
}

// 处理连接
func (s *Server) handleConn(conn net.Conn) {
	// 将连接包装成缓冲处理器，提高读取性能
	reader := bufio.NewReader(conn)
	defer conn.Close()

	for {
		command, args, err := readRequestFrom(reader)
		if err != nil {
			if err == ProtocolVersionMismatchErr {
				continue
			}
			return
		}

		// 处理请求
		reply, body, err := s.handleRequest(command, args)
		if err != nil {
			writeErrorResponseTo(conn, err.Error())
			continue
		}

		// 发送处理结果响应
		_, err = writeResponseTo(conn, reply, body)
		if err != nil {
			continue
		}
	}
}

func (s *Server) handleRequest(command byte, args [][]byte) (reply byte, body []byte, err error) {
	// 从命令集合中选出对应的处理器
	handle, ok := s.handlers[command]
	if !ok {
		return ErrorReply, nil, commandHandlerNotFoundErr
	}

	// 将处理结果返回
	body, err = handle(args)
	if err != nil {
		return ErrorReply, body, err
	}
	return SuccessReply, body, err
}

func (s *Server) Close() error {
	if s.listener == nil {
		return nil
	}
	return s.listener.Close()
}

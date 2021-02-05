package vex

import "errors"

const (
	ProtocolVersion        = byte(1) // 协议版本号
	headerLengthInProtocol = 6       // 头部占用字节
	argsLengthInProtocol   = 4       // 参数个数占用字节
	argLengthInProtocol    = 4       // 协议中参数长度占用字节
	bodyLengthInProtocol   = 4       // 协议体长度占用字节数
)

var (
	ProtocolVersionMismatchErr = errors.New("protocol version between client and server doesn't match")
)

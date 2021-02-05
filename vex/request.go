package vex

import (
	"encoding/binary"
	"io"
)

func readRequestFrom(reader io.Reader) (command byte, args [][]byte, err error) {
	// 读取头部
	header := make([]byte, headerLengthInProtocol)
	// ReadFull 方法，如果数据没有读满，会等待
	_, err = io.ReadFull(reader, header)
	if err != nil {
		return 0, nil, err
	}
	// 头部第一个字节是协议版本号，取出来判断协议版本号是否一致
	version := header[0]
	if version != ProtocolVersion {
		return 0, nil, ProtocolVersionMismatchErr
	}

	// 头部的第二个字节是命令 ,后面四个字节是参数
	command = header[1]
	header = header[2:]

	// 所有整数到字节数组的转换使用大端形式，所以这里使用 BigEndian 将头部后四个字节转换成一个uint32 数字
	// argsLength: 参数的个数
	argsLength := binary.BigEndian.Uint32(header)
	args = make([][]byte, argsLength)
	if argsLength > 0 {
		// 读取参数长度，同样使用大端形式处理
		argLength := make([]byte, argLengthInProtocol)
		for i := uint32(0); i < argsLength; i++ {
			_, err = io.ReadFull(reader, argLength)
			if err != nil {
				return 0, nil, err
			}
			arg := make([]byte, binary.BigEndian.Uint32(argLength))
			_, err = io.ReadFull(reader, arg)
			if err != nil {
				return 0, nil, err
			}
			args[i] = arg
		}
	}
	return command, args, nil
}

func writeRequestTo(writer io.Writer, command byte, args [][]byte) (int, error) {
	request := make([]byte, headerLengthInProtocol)
	request[0] = ProtocolVersion
	request[1] = command
	binary.BigEndian.PutUint32(request[2:], uint32(len(args)))

	if len(args) > 0 {
		// 将参数都添加到缓冲区
		argLength := make([]byte, argLengthInProtocol)
		for _, arg := range args {
			binary.BigEndian.PutUint32(argLength, uint32(len(arg)))
			request = append(request, argLength...)
			request = append(request, arg...)
		}
	}
	return writer.Write(request)
}

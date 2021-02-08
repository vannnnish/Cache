package vex

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	SuccessReply = 0
	ErrorReply   = 1
)

func readResponseFrom(reader io.Reader) (reply byte, body []byte, err error) {
	// 读取指定字节数据
	header := make([]byte, headerLengthInProtocol)
	_, err = io.ReadFull(reader, header)
	if err != nil {
		return ErrorReply, nil, err
	}

	version := header[0]
	if version != ProtocolVersion {
		return ErrorReply, nil, errors.New("response " + ProtocolVersionMismatchErr.Error())
	}

	// reply: 命令
	reply = header[1]
	// 响应体长度
	header = header[2:]
	body = make([]byte, binary.BigEndian.Uint32(header))
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return ErrorReply, nil, err
	}

	return reply, body, nil
}

// 将响应写入到writer
func writeResponseTo(writer io.Writer, reply byte, body []byte) (int, error) {
	bodyLengthBytes := make([]byte, bodyLengthInProtocol)
	binary.BigEndian.PutUint32(bodyLengthBytes, uint32(len(body)))

	response := make([]byte, 2, headerLengthInProtocol+len(body))
	response[0] = ProtocolVersion
	response[1] = reply
	response = append(response, bodyLengthBytes...)
	response = append(response, body...)
	return writer.Write(response)
}

// 向writer 写入错误信息为msg 的响应
func writeErrorResponseTo(writer io.Writer, msg string) (int, error) {
	return writeResponseTo(writer, ErrorReply, []byte(msg))
}

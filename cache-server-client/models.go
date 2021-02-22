package cache_server_client

import (
	"encoding/json"
)

type Status struct {
	Count     int   `json:"count"`
	KeySize   int64 `json:"keySize"`
	ValueSize int64 `json:"valueSize"`
}

// request 请求结构体
type request struct {
	// 命令
	command byte
	// 执行的参数
	args [][]byte
	// 用于接收结果的管道
	resultChan chan *Response
}

type Response struct {
	// 响应体
	Body []byte
	// 响应的错误
	Err error
}

func (r *Response) ToStatus() (*Status, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	status := &Status{}
	return status, json.Unmarshal(r.Body, status)
}

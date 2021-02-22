package services

import (
	"cache/caches"
	"cache/vex"
	"encoding/binary"
	"encoding/json"
)

type TCPClient struct {
	client *vex.Client
}

func NewTCPClient(address string) (*TCPClient, error) {
	client, err := vex.NewClient("tcp", address)
	if err != nil {
		return nil, err
	}
	return &TCPClient{client: client}, nil
}

func (tc *TCPClient) Get(key string) ([]byte, error) {
	return tc.client.Do(getCommand, [][]byte{[]byte(key)})
}
func (tc *TCPClient) Set(key string, value []byte, ttl int64) error {
	ttlBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(ttlBytes, uint64(ttl))
	_, err := tc.client.Do(setCommand, [][]byte{ttlBytes, []byte(key), value})
	return err
}

func (tc *TCPClient) Delete(key string) error {
	_, err := tc.client.Do(deleteCommand, [][]byte{[]byte(key)})
	return err
}

func (tc *TCPClient) Status() (*caches.Status, error) {
	body, err := tc.client.Do(statusCommand, nil)
	if err != nil {
		return nil, err
	}
	status := caches.NewStatus()
	err = json.Unmarshal(body, status)
	return status, err
}

func (tc *TCPClient) Close() error {
	return tc.client.Close()
}

func (tc *TCPClient) Nodes() ([]string, error) {
	body, err := tc.client.Do(nodesCommand, nil)
	if err != nil {
		return nil, err
	}
	var nodes []string
	err = json.Unmarshal(body, &nodes)
	return nodes, err

}

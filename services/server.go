package services

import "cache/caches"

const (
	APIVersion = "v1"
)

type Server interface {
	Run(address string) error
}

func NewServer(serverType string, cache *caches.Cache) Server {
	if serverType == "tcp" {
		return NewTcpServer(cache)
	}
	return NewHTTPServer(cache)
}


package services

import "cache/caches"

const (
	APIVersion = "v1"
)

type Server interface {
	Run() error
}

func NewServer(cache *caches.Cache, options Options) (Server, error) {
	if options.ServerType == "tcp" {
		return NewTcpServer(cache, &options)
	}
	return NewHTTPServer(cache, &options)
}

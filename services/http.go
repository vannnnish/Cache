package services

import (
	"cache/caches"
	"cache/helpers"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
)

type HTTPServer struct {
	*node
	cache   *caches.Cache
	options *Options
}

func NewHTTPServer(cache *caches.Cache, options *Options) (*HTTPServer, error) {
	n, err := newNode(options)
	if err != nil {
		return nil, err
	}
	return &HTTPServer{
		node:    n,
		cache:   cache,
		options: options,
	}, nil
}

func (hs *HTTPServer) Run() error {
	return http.ListenAndServe(helpers.JoinAddressAndPort(hs.options.Address, hs.options.Port), hs.routerHandler())
}

func wrapUriWithVersion(uri string) string {
	return path.Join("/", APIVersion, uri)
}

func (hs *HTTPServer) routerHandler() http.Handler {
	router := httprouter.New()
	router.GET(wrapUriWithVersion("/cache/:key"), hs.getHandler)
	router.PUT(wrapUriWithVersion("/cache/:key"), hs.setHandler)
	router.DELETE(wrapUriWithVersion("/cache/:key"), hs.deleteHandler)
	router.GET(wrapUriWithVersion("/status"), hs.statusHandler)

	router.GET(wrapUriWithVersion("/nodes"), hs.nodesHandler)
	return router
}

func (hs *HTTPServer) getHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	node, err := hs.selectNode(key)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 判断这个 key 所属的物理节点是否是当前节点， 如果不是， 需要响应重定向信息给客户端，并告知正确的节点地址
	if !hs.isCurrentNode(node) {
		writer.Header().Set("Location", node+request.RequestURI)
		writer.WriteHeader(http.StatusTemporaryRedirect)
	}

	// 当前节点处理
	value, ok := hs.cache.Get(key)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.Write(value)
}

func (hs *HTTPServer) setHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	// 使用一致性哈希选择出key所在的物理节点
	key := params.ByName("key")
	node, err := hs.selectNode(key)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 判断这个key所属的是否是当前节点，如果不是，需要重新定向给客户端，并告知正确的节点地址
	if !hs.isCurrentNode(node) {
		writer.Header().Set("Location", node+request.RequestURI)
		writer.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	// 当前节点处理
	value, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	ttl, err := ttlOf(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = hs.cache.SetWithTTL(key, value, ttl)
	if err != nil {
		writer.WriteHeader(http.StatusRequestEntityTooLarge)
		writer.Write([]byte("Error:" + err.Error()))
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

func ttlOf(request *http.Request) (int64, error) {
	ttls, ok := request.Header["Ttl"]
	if !ok || len(ttls) < 1 {
		return caches.NeverDie, nil
	}
	return strconv.ParseInt(ttls[0], 10, 64)
}

func (hs *HTTPServer) deleteHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	node, err := hs.selectNode(key)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !hs.isCurrentNode(node) {
		writer.Header().Set("Location", node+request.RequestURI)
		writer.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	// 当前节点处理
	err = hs.cache.Delete(key)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (hs *HTTPServer) statusHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	status, err := json.Marshal(hs.cache.Status())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(status)
}

func (hs *HTTPServer) nodesHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	nodes, err := json.Marshal(hs.nodes())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(nodes)
}

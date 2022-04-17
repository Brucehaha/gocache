package servers

import (
	"cache-server/caches"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HTTPServer struct {
	cache *caches.Cache
}

func NewHTTPServer(cache *caches.Cache) *HTTPServer {
	return &HTTPServer{
		cache: cache,
	}
}

func (hs *HTTPServer) Run(address string) error {
	return http.ListenAndServe(address, hs.routerHander())

}
func (hs *HTTPServer) routerHander() http.Handler {
	router := httprouter.New()
	router.GET("/cache/:key", hs.getHandler)
	router.PUT("/cache/:key", hs.setHandler)
	router.DELETE("/cache/:key", hs.deleteHandler)
	router.GET("/status", hs.statusHandler)
	return router
}

func (hs *HTTPServer) getHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	value, ok := hs.cache.Get(key)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
	}
	writer.Write(value)

}
func (hs *HTTPServer) setHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	value, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	hs.cache.Set(key, value)
}

func (hs *HTTPServer) deleteHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	hs.cache.Delete(key)

}
func (hs *HTTPServer) statusHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	status, err := json.Marshal(map[string]interface{}{
		"count": hs.cache.Count(),
	})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(status)
}

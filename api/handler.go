package api

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/jonfriesen/redis-proxy/cache"
)

type handler struct {
	Cache *cache.Cache
}

var (
	ErrNotFound               = errors.New("Error: key-value pair not found")
	ErrRquestTypeNotSupported = errors.New("Error: HTTP Request type not supported")
)

func New(c *cache.Cache) http.Handler {
	mux := http.NewServeMux()

	h := handler{
		Cache: c,
	}

	mux.Handle("/v1/get/", wrapper(h.get))

	return mux
}

func (h *handler) get(w io.Writer, r *http.Request) (interface{}, int, error) {
	switch r.Method {
	case "GET":
		rKey := strings.TrimPrefix(r.URL.Path, "/v1/get/")
		log.Printf("Looking up %v", rKey)

		v, err := h.Cache.Get(rKey)
		if err != nil {
			return nil, http.StatusNotFound, ErrNotFound
		}

		return v, http.StatusOK, nil
	}

	return nil, http.StatusBadRequest, ErrRquestTypeNotSupported
}

package api

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/jonfriesen/redis-proxy/cache"
	"github.com/jonfriesen/redis-proxy/storage"
)

type handler struct {
	DataSource *storage.Storage
	Cache      *cache.Cache
}

var (
	ErrNotFound               = errors.New("Error: key-value pair not found")
	ErrRquestTypeNotSupported = errors.New("Error: HTTP Request type not supported")
)

func New(ds *storage.Storage, c *cache.Cache) http.Handler {
	mux := http.NewServeMux()

	h := handler{
		DataSource: ds,
		Cache:      c,
	}

	mux.Handle("/v1/get/", wrapper(h.get))

	return mux
}

func (h *handler) get(w io.Writer, r *http.Request) (interface{}, int, error) {
	switch r.Method {
	case "GET":
		rKey := strings.TrimPrefix(r.URL.Path, "/v1/get/")
		log.Printf("Looking up %v", rKey)

		h.Cache.Lock()
		defer h.Cache.Unlock()
		v, err := h.Cache.Get(rKey)
		if err == cache.ErrNotFound {
			log.Printf("Lookup not in cache %v", rKey)

			v, err = (*h.DataSource).Get(rKey)
			if err == storage.ErrNotFound {
				log.Printf("Lookup not in storage %v", rKey)
				return nil, http.StatusNotFound, ErrNotFound
			}

			if v != "" {
				log.Println("Pushing key-value pair into cache")
				err = h.Cache.Push(rKey, v)
			}
		}
		if err == cache.ErrUnlocked {
			log.Fatalf("Critical Error: %v", err)
		}

		return v, http.StatusOK, nil
	}

	return nil, http.StatusBadRequest, ErrRquestTypeNotSupported
}

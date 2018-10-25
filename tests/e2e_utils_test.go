package tests

import (
	"strings"

	"github.com/jonfriesen/redis-proxy/storage"
)

type storageMiddleware struct {
	table    map[string]int
	callList []string
	actual   *storage.Storage
}

func NewStorageMiddleware(s *storage.Storage) *storageMiddleware {

	i := storageMiddleware{
		table:    make(map[string]int),
		callList: []string{},
		actual:   s,
	}

	return &i
}

func (i *storageMiddleware) Get(key string) (string, error) {
	i.catch([]string{"GET", key})

	s, err := (*i.actual).Get(key)
	return s, err

}

func (i *storageMiddleware) Put(key, value string) error {
	i.catch([]string{"PUT", key, value})

	err := (*i.actual).Put(key, value)

	return err
}

func (i *storageMiddleware) catch(p []string) {
	action := strings.Join(p, "$")
	_, e := i.table[action]
	if !e {
		i.table[action] = 1
	} else {
		i.table[action]++
	}

	i.callList = append(i.callList, action)
}

func (i *storageMiddleware) asStorage() *storage.Storage {
	var wrapper storage.Storage
	wrapper = i

	return &wrapper
}

func seedDatasource(m map[string]string, ds *storage.Storage) {
	for k, v := range m {
		(*ds).Put(k, v)
	}
}

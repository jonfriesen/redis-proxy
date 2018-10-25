package inmem

import "github.com/jonfriesen/redis-proxy/storage"

type inmem struct {
	table map[string]string
}

func New(seed map[string]string) *storage.Storage {

	i := inmem{}

	if seed != nil {
		i = inmem{seed}
	}

	var s storage.Storage
	s = &i

	return &s
}

func (i *inmem) Get(key string) (string, error) {
	v, exists := i.table[key]
	if !exists {
		return "", storage.ErrNotFound
	}

	return v, nil
}

func (i *inmem) Put(key, value string) error {
	i.table[key] = value
	return nil
}

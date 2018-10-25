package redis

import (
	"fmt"
	"log"

	"github.com/jonfriesen/redis-proxy/storage"
	"github.com/mediocregopher/radix.v2/pool"
)

type redisPool struct {
	db *pool.Pool
}

func New(host, port string) (*storage.Storage, error) {
	db, err := pool.New("tcp", fmt.Sprintf("%v:%v", host, port), 10)
	if err != nil {
		log.Fatalf("Failed to create Redis pool: %v", err)
		return nil, err
	}

	rp := redisPool{db}

	var s storage.Storage
	s = &rp
	return &s, nil
}

func (r *redisPool) Get(key string) (string, error) {
	c, err := r.db.Get()
	if err != nil {
		return "", err
	}
	defer r.db.Put(c) // return connection to the pool

	v, err := c.Cmd("GET", key).Str()
	if err != nil && err.Error() != "No Record Found" {
		return "", storage.ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return v, nil
}

func (r *redisPool) Put(key, value string) error {
	c, err := r.db.Get()
	if err != nil {
		return err
	}
	defer r.db.Put(c) // return connection to the pool

	_, err = c.Cmd("SET", key, value).Str()
	if err != nil {
		return err
	}

	return nil
}

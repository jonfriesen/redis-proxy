package storage

import "errors"

// Storage represents a generic interface as a contract
// to substitute data sources. For example a Memcache,
// MongoDB, or other key-value store implementation
// could be build here.
type Storage interface {
	Get(string) (string, error)
	Put(string, string) error
}

var ErrNotFound = errors.New("Error: Could not find record")

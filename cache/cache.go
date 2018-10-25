package cache

import (
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("No Record Found")

type record struct {
	key    string
	value  string
	expiry time.Time
}

type node struct {
	record *record
	parent *node
	child  *node
}

type Cache struct {
	table   map[string]*node
	head    *node
	tail    *node
	maxKeys int32
	maxAge  time.Duration
	mtx     *sync.Mutex
}

func New(mkeys int32, mage time.Duration) *Cache {
	return &Cache{
		make(map[string]*node),
		nil,
		nil,
		mkeys,
		mage,
		&sync.Mutex{},
	}
}

func (c *Cache) Push(key, value string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	n := &node{
		record: &record{
			key:    key,
			value:  value,
			expiry: time.Now().Add(c.maxAge),
		},
	}

	c.table[n.record.key] = n

	c.add(n)

	// clean up if necessary
	if int32(len(c.table)) > c.maxKeys {
		d := c.head
		if d != nil {
			c.evict(d)
			delete(c.table, d.record.key)
		}
	}
}

func (c *Cache) Get(key string) (string, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	n, exists := c.table[key]
	if !exists {
		return "", ErrNotFound
	}
	if time.Now().After(n.record.expiry) {
		return "", ErrNotFound
	}

	// could defer these til after the return
	c.evict(n)
	c.add(n)

	return n.record.value, nil
}

// add will add a node to the list Cache
func (c *Cache) add(n *node) {

	// prep node for insertion
	n.parent = nil
	n.child = c.tail

	// redirect current tail to point to new tail
	if n.child != nil {
		n.child.parent = n
	}

	c.tail = n

	// handle empty Cache
	if c.head == nil {
		c.head = n
	}

}

// evict will remove a node from the list
func (c *Cache) evict(n *node) {

	//check if head/tail
	if n == c.tail {
		c.tail = n.child
	}
	if n == c.head {
		c.head = n.parent
	}

	// remap nodes on either side
	if n.parent != nil {
		n.parent.child = n.child
	}
	if n.child != nil {
		n.child.parent = n.parent
	}

}

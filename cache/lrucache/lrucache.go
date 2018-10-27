package lrucache

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("No Record Found")
	ErrUnlocked = errors.New("LRUCache much be locked with LRUCache.Lock()")
)

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

type LRUCache struct {
	table   map[string]*node
	head    *node
	tail    *node
	maxKeys int32
	maxAge  time.Duration
	mtx     *sync.Mutex
	locked  bool
}

func New(mkeys int32, mage time.Duration) *LRUCache {
	return &LRUCache{
		make(map[string]*node),
		nil,
		nil,
		mkeys,
		mage,
		&sync.Mutex{},
		false,
	}
}

func (c *LRUCache) Push(key, value string) error {
	if !c.locked {
		return ErrUnlocked
	}

	// remove duplicates from linklist
	oldNode, exists := c.table[key]
	if exists {
		c.evict(oldNode)
	}

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

	return nil
}

func (c *LRUCache) Get(key string) (string, error) {
	if !c.locked {
		return "", ErrUnlocked
	}

	n, exists := c.table[key]
	if !exists {
		return "", ErrNotFound
	}
	if time.Now().After(n.record.expiry) {
		c.evict(n)
		return "", ErrNotFound
	}

	// could defer these til after the return
	c.evict(n)
	c.add(n)

	return n.record.value, nil
}

func (c *LRUCache) Lock() {
	c.mtx.Lock()
	c.locked = true
}

func (c *LRUCache) Unlock() {
	c.mtx.Unlock()
	c.locked = false
}

// add will add a node to the list Cache
func (c *LRUCache) add(n *node) {

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
func (c *LRUCache) evict(n *node) {

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

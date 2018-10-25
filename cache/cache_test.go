package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCache_add(t *testing.T) {
	type fields struct {
		table   map[string]*node
		head    *node
		tail    *node
		maxKeys int32
		maxAge  time.Duration
		mtx     *sync.Mutex
	}
	type args struct {
		n *node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "add node to empty cache",
			fields: fields{
				table:   make(map[string]*node),
				head:    nil,
				tail:    nil,
				maxKeys: 5,
				maxAge:  1 * time.Minute,
				mtx:     &sync.Mutex{},
			},
			args: args{
				n: &node{
					record: &record{
						key:    "testkey",
						value:  "testvalue",
						expiry: time.Now().Add(1 * time.Millisecond),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				table:   tt.fields.table,
				head:    tt.fields.head,
				tail:    tt.fields.tail,
				maxKeys: tt.fields.maxKeys,
				maxAge:  tt.fields.maxAge,
				mtx:     tt.fields.mtx,
			}
			c.add(tt.args.n)
			fmt.Printf("%+v", c.tail.record.value)
			if c.tail != tt.args.n && c.head != tt.args.n {
				t.Error("Record not present in linked list")
			}
		})
	}
}

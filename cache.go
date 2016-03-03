package kvcache

import (
	// "fmt"
	"sync"
	"time"
)

var (
	lock   sync.Mutex
	tables map[string]*Table
)

func Cache(name string) *Table {
	t, ok := tables[name]
	if !ok {
		t = &Table{
			items:          make(map[interface{}]*Item),
			interval:       0 * time.Second,
			addCallback:    nil,
			deleteCallback: nil,
			log:            nil,
		}
	}
	return t
}

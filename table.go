package kvcache

import (
	"log"
	"sync"
	"time"
)

type Table struct {
	sync.Mutex
	items          map[interface{}]*Item
	interval       time.Duration
	log            *log.Logger
	addCallback    func(interface{})
	deleteCallback func(interface{})
}

func (t *Table) Count() int {
	t.Lock()
	defer t.Unlock()
	return len(t.items)
}

func (t *Table) SetAddCallback(f func(interface{})) {
	t.Lock()
	defer t.Unlock()
	t.addCallback = f
}

func (t *Table) SetDeleteCallback(f func(interface{})) {
	t.Lock()
	defer t.Unlock()
	t.deleteCallback = f
}

func (t *Table) Add(key interface{}, value interface{}, duration time.Duration) *Item {
	t.Lock()
	item := CreateItem(key, value, duration)
	t.items[key] = &item
	t.Unlock()
	if t.addCallback != nil {
		t.addCallback(item)
	}
	if duration > 0 && (t.interval == 0 || duration < t.interval) {
		t.expireCheck()
	}
	return &item
}

func (t *Table) Get(key interface{}) (interface{}, error) {
	t.Lock()
	defer t.Unlock()
	v, ok := t.items[key]
	if !ok {
		return nil, NotFoundedKey
	}
	v.accessCount++
	return v.data, nil
}

func (t *Table) Exist(key interface{}) bool {
	t.Lock()
	defer t.Unlock()
	_, ok := t.items[key]
	if !ok {
		return false
	}
	return true
}

func (t *Table) Flush() {
	t.Lock()
	defer t.Unlock()
	t.items = make(map[interface{}]*Item)
	t.interval = 0
	t.writeLog("Flush items")
}

func (t *Table) Delete(key interface{}) (*Item, error) {
	t.Lock()
	defer t.Unlock()
	i, ok := t.items[key]
	if !ok {
		return nil, NotFoundedKey
	}
	if t.deleteCallback != nil {
		t.deleteCallback(i)
	}
	if i.deleteCallback != nil {
		i.deleteCallback(i)
	}
	delete(t.items, key)
	t.writeLog("Delete item key:", key, "Create on:", i.createTime, "Data:", i.data)
	return i, nil

}

func (t *Table) expireCheck() {
	t.Lock()
	items := t.items
	t.Unlock()
	now := time.Now()
	smallestTime := t.interval
	for key, item := range items {
		item.Lock()
		duration := item.duration
		lastTime := item.lastTime
		item.Unlock()
		if duration > 0 && now.Sub(lastTime) > duration {
			t.Delete(key)
		}
		if duration > 0 && duration < smallestTime {
			smallestTime = duration
		}
	}
	t.Lock()
	if smallestTime < t.interval {
		t.interval = smallestTime
	}
	time.AfterFunc(smallestTime, func() {
		go t.expireCheck()
	})
	t.Unlock()
}

func (t *Table) SetLogger(l *log.Logger) {
	t.Lock()
	defer t.Unlock()
	t.log = l
}

func (t *Table) writeLog(v ...interface{}) {
	if t.log == nil {
		return
	}
	t.log.Println(v)
}

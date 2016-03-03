package kvcache

import (
	"sync"
	"time"
)

type Item struct {
	sync.Mutex
	key            interface{}
	data           interface{}
	duration       time.Duration
	createTime     time.Time
	lastTime       time.Time
	accessCount    int64
	deleteCallback func(interface{})
}

func CreateItem(key interface{}, value interface{}, duration time.Duration) Item {
	t := time.Now()
	return Item{
		key:        key,
		data:       value,
		duration:   duration,
		createTime: t,
		lastTime:   t,
	}
}

func (i *Item) KeepAlive() {
	i.Lock()
	defer i.Unlock()
	i.accessCount++
	i.lastTime = time.Now()
}

func (i *Item) LastTime() time.Time {
	i.Lock()
	defer i.Unlock()
	return i.lastTime
}

func (i *Item) CreateTime() time.Time {
	return i.createTime
}

func (i *Item) AccessCount() int64 {
	i.Lock()
	defer i.Unlock()
	return i.accessCount
}

func (i *Item) SetDeleteCallback(f func(interface{})) {
	i.Lock()
	defer i.Unlock()
	i.deleteCallback = f
}

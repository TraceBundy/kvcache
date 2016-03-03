package main

import (
	"fmt"
	"github.com/TraceBundy/kvcache"
	"time"
)

func main() {
	t := kvcache.Cache("hello")
	t.SetAddCallback(func(i interface{}) {
		fmt.Println("add new item", i)
	})
	t.SetDeleteCallback(func(i interface{}) {
		fmt.Println("delete item", i)
	})
	t.Add("hello", 1, 1*time.Second)
	v, ok := t.Get("hello")
	if ok == nil {
		fmt.Println(v)
	} else {
		fmt.Printf("%s\n", ok)
	}
	time.Sleep(6000 * time.Second)
}

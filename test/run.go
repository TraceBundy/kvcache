package main

import (
	"fmt"
	"github.com/TraceBundy/kvcache"
	"time"
)

func main() {
	t := kvcache.Cache("hello")
	t.Add("hello", 1, -1)
	v, ok := t.Get("hello")
	if ok == nil {
		fmt.Println(v)
	} else {
		fmt.Printf("%s\n", ok)
	}
	t.Delete("hello")
	time.Sleep(6 * time.Second)
}

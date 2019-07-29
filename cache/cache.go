package cache

import (
	"errors"
	"fmt"
	"github.com/mitesoro/library/stat/prom"
	"runtime"
	"sync"
)

var (
	// ErrFull cache internal chan full.
	ErrFull = errors.New("cache chan full")
	stats   = prom.BusinessInfoCount
)

// Cache 通过chan缓存异步保存数据。
type Cache struct {
	ch     chan func()
	worker int
	waiter sync.WaitGroup
}

// New 新建
func New(worker, size int) *Cache {
	if worker <= 0 {
		worker = 1
	}
	c := &Cache{
		ch:     make(chan func(), size),
		worker: worker,
	}
	c.waiter.Add(worker)
	for i := 0; i < worker; i++ {
		go c.proc()
	}
	return c
}

func (c *Cache) proc() {
	defer c.waiter.Done()
	for {
		f := <-c.ch
		if f == nil {
			return
		}
		wrapFunc(f)()
		stats.State("cache_channel", int64(len(c.ch)))
	}
}

//wrapFunc 包装方法 panic处理
func wrapFunc(f func()) (res func()) {
	res = func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64*1024)
				buf = buf[:runtime.Stack(buf, false)]
				//log.Error("panic in cache proc,err: %s,stack: %s", f, buf)
				fmt.Println(string(buf))
			}
		}()
		f()
	}
	return
}

// Save 保存
func (c *Cache) Save(f func()) (err error) {
	if f == nil {
		return
	}
	select {
	case c.ch <- f:
	default:
		err = ErrFull
	}
	stats.State("cache_channel", int64(len(c.ch)))
	return
}

// Close 关闭
func (c *Cache) Close() (err error) {
	for i := 0; i < c.worker; i++ {
		c.ch <- nil
	}
	c.waiter.Wait()
	return
}

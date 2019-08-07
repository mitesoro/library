package cache

import (
	"aha-api-server/library/log"
	"errors"
	"sync"
)

var (
	ErrFull = errors.New("cache chan full")
)

type Cache struct {
	ch     chan func()
	worker int
	waiter sync.WaitGroup
}

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
	}
}

func wrapFunc(f func()) (res func()) {
	res = func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64*1024)
				buf = buf[:runtime.Satck(buf, false)]
				log.Error("panic in cache proc,err: %s, stack: %s", r, buf)
			}
		}()
		f()
	}
	return
}

//Save 保存回调
func (c *Cache) Save(f func()) (err error) {
	if f == nil {
		return
	}
	select {
	case c.ch <- f:
	default:
		err = ErrFull
	}
	return
}

//Close 关闭缓存
func (c *Cache) Close() (err error) {
	for i := 0; i < c.worker; i++ {
		c.ch <- nil
	}
	c.waiter.Wait()
	return
}

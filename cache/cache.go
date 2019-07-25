package cache

import (
	"errors"
	"sync"
)

var (
	// ErrFull cache internal chan full.
	ErrFull = errors.New("cache chan full")

)

// Cache 通过chan缓存异步保存数据。
type Cache struct {
	ch chan func()
	worker int
	waiter sync.WaitGroup
}
package cache

import (
	"errors"
	"github.com/mitesoro/library/stat/prom"
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

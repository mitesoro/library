package log

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileLog(t *testing.T) {
	Init(nil) // debug flag: log.dir={path}
	//initFile()
	Info("hello %s", "test")
	assert.Equal(t, nil, Close())
}

func initStdout() {
	conf := &Config{
		Stdout: true,
	}
	Init(conf)
}

func initFile() {
	conf := &Config{
		Dir: "/Users/liuhaigang/go/src/github.com/mitesoro/logs",
		// VLevel:  2,
		Module: map[string]int32{"log_test": 1},
	}
	Init(conf)
}

type TestLog struct {
	A string
	B int
	C string
	D string
}

func testLog(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		Error("hello %s", "world")
		Errorv(context.Background(), KV("key", 2222222), KV("test2", "test"))
	})
	t.Run("Warn", func(t *testing.T) {
		Warn("hello %s", "world")
		Warnv(context.Background(), KV("key", 2222222), KV("test2", "test"))
	})
	t.Run("Info", func(t *testing.T) {
		Info("hello %s", "world")
		Infov(context.Background(), KV("key", 2222222), KV("test2", "test"))
	})
}

func TestFile(t *testing.T) {
	initFile()
	//testLog(t)
	Info("hello %s", "test")
	//Infov(context.Background(), KV("key", 2222222), KV("test2", "test"))
	assert.Equal(t, nil, Close())
}

func TestStdout(t *testing.T) {
	initStdout()
	testLog(t)
	assert.Equal(t, nil, Close())
}

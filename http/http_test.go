package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestHttp(t *testing.T) {
	engine := gin.New()
	NewService(engine, "get", "/ping", Ping)
	engine.Run(":1234")
}

func Ping(ctx *Context) {
	err := errors.New("错误测试")
	ctx.Error(err)
}

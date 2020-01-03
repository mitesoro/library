package log

import (
	"fmt"
	"github.com/mitesoro/library/http"
)

func Logger() http.HandlerFunc {

	return func(c *http.Context) {
		fmt.Println("latency")
	}
}

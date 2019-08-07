package log

import (
	"aha-api-server/library/http"
	"fmt"
)

func Logger() http.HandlerFunc {

	return func(c *http.Context) {
		fmt.Println("latency")
	}
}

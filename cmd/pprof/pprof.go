package pprof

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func Listen(port int) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	return http.ListenAndServe(addr, nil)
}

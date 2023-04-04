package testing

import (
	"net/http"
)

func HelloHandler(resp http.ResponseWriter, req *http.Request) {
	_, _ = resp.Write([]byte("hello world"))
}

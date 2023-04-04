package std

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"
)

type HTTPProxy struct {
	proxy *httputil.ReverseProxy
}

func NewHTTPProxy(target string) (*HTTPProxy, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	return &HTTPProxy{httputil.NewSingleHostReverseProxy(u)}, nil
}

func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.proxy.ServeHTTP(w, req)
}

func TestHTTPProxy(t *testing.T) {
	proxy, err := NewHTTPProxy("http:127.0.0.1:8080")
	if err != nil {
		return
	}
	http.Handle("/", proxy)
	_ = http.ListenAndServe(":8081", nil)
}

package proxy

import (
	"Golearn/src/projects/tinybalancer/balancer"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var (
	XRealIP       = http.CanonicalHeaderKey("X-Real-IP")
	XProxy        = http.CanonicalHeaderKey("X-Proxy")
	XForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
)

var (
	ReverseProxy = "Balancer-Reverse-Proxy"
)

type HTTPProxy struct {
	hostMap      map[string]*httputil.ReverseProxy // 主机对反向代理的映射
	lb           balancer.Balancer                 // 负载均衡器
	sync.RWMutex                                   // 读写锁，保护alive的并发安全
	alive        map[string]bool                   // 表示反向代理的主机是否处于健康状态
}

// NewHTTPProxy targetHosts代理主机，algorithm负载均衡算法
func NewHTTPProxy(targetHosts []string, algorithm string) (*HTTPProxy, error) {
	httpProxy := &HTTPProxy{
		hostMap: make(map[string]*httputil.ReverseProxy),
		alive:   make(map[string]bool),
	}
	hosts := make([]string, 0)

	for _, targetHost := range targetHosts {
		targetUrl, err := url.Parse(targetHost)
		if err != nil {
			return nil, err
		}

		// 单主机反向代理
		proxy := httputil.NewSingleHostReverseProxy(targetUrl)
		originDirector := proxy.Director
		// 修改请求
		proxy.Director = func(req *http.Request) {
			originDirector(req)
			req.Header.Set(XProxy, ReverseProxy)
			req.Header.Set(XRealIP, GetIP(req))
		}

		host := GetHost(targetUrl)
		httpProxy.alive[host] = true
		httpProxy.hostMap[host] = proxy
		hosts = append(hosts, host)
	}

	// 根据算法构建负载均衡器
	lb, err := balancer.Build(algorithm, hosts)
	if err != nil {
		return nil, err
	}
	httpProxy.lb = lb

	return httpProxy, nil
}

func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("proxy causes panic: %s", err)
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(err.(error).Error()))
		}
	}()

	// 根据IP地址选择合适的主机
	host, err := h.lb.Balance(GetIP(r))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(fmt.Sprintf("balance error: %s", err.Error())))
		return
	}

	h.lb.Inc(host)
	defer h.lb.Done(host)
	h.hostMap[host].ServeHTTP(w, r)
}

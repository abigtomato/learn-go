package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ConnectionTimeout = 5 * time.Second

// GetIP
// 若客户端IP 为 192.168.1.1 通过代理 192.168.2.5 和 192.168.2.6
// X-Forwarded-For的值可能为 [192.168.2.5, 192.168.2.6]
// X-Real-IP的值为 192.168.1.1
func GetIP(req *http.Request) string {
	clientIP, _, _ := net.SplitHostPort(req.RemoteAddr)
	if len(req.Header.Get(XForwardedFor)) != 0 {
		// 尝试从X-Forwarded-For获取客户端IP地址
		xff := req.Header.Get(XForwardedFor)
		s := strings.Index(xff, ", ")
		if s == -1 {
			s = len(req.Header.Get(XForwardedFor))
		}
		clientIP = xff[:s]
	} else if len(req.Header.Get(XRealIP)) != 0 {
		// 尝试从X-Real-IP获取客户端IP地址
		clientIP = req.Header.Get(XRealIP)
	}
	return clientIP
}

func GetHost(url *url.URL) string {
	if _, _, err := net.SplitHostPort(url.Host); err != nil {
		return url.Host
	}
	if url.Scheme == "http" {
		return fmt.Sprintf("%s:%s", url.Host, "80")
	} else if url.Scheme == "https" {
		return fmt.Sprintf("%s:%s", url.Host, "443")
	}
	return url.Host
}

// IsBackendAlive 判断后台连接的活性
func IsBackendAlive(host string) bool {
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return false
	}
	resolve := fmt.Sprintf("%s:%d", addr.IP, addr.Port)
	conn, err := net.DialTimeout("tcp", resolve, ConnectionTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

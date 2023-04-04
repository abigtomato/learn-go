package geecache

import (
	"Golearn/src/projects/geecache/consistenthash"
	"Golearn/src/projects/geecache/geecachepb"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

// HTTPPool 分布式缓存系统节点通信 具备提供http服务和客户端的能力
type HTTPPool struct {
	self        string                 // 记录自己的地址
	basePath    string                 // 节点间通信的前缀
	mu          sync.Mutex             // 互斥锁
	peers       *consistenthash.Map    // 节点哈希环
	httpGetters map[string]*httpGetter // 远程节点与其http获取器的映射
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...any) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// 提供http服务的能力
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 前缀判断
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	p.Log("%s %s", r.Method, r.URL.Path)

	// 路径截取 /<basePath>/<groupName>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	// 缓存组
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	// 缓存值
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// protobuf编码
	body, err := proto.Marshal(&geecachepb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err = w.Write(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Set 设置远程peer节点
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 节点上环
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)

	// 节点绑定http客户端
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// PickPeer 基于http客户端的peer节点选择器
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 从哈希环上选择
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)

// http客户端
type httpGetter struct {
	baseURL string // 需要访问的远程节点地址
}

// Get 发送http get请求从其他peer节点获取缓存
func (h *httpGetter) Get(in *geecachepb.Request, out *geecachepb.Response) error {
	u := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(in.GetGroup()), url.QueryEscape(in.GetKey()))

	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	// protobuf解码
	if err := proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}

	return nil
}

var _ PeerGetter = (*httpGetter)(nil)

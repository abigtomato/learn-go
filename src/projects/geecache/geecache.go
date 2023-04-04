package geecache

import (
	"Golearn/src/projects/geecache/geecachepb"
	"Golearn/src/projects/geecache/singleflight"
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 回调函数（缓存不存在时，调用用户自定义函数，获得源数据）
type GetterFunc func(key string) ([]byte, error)

// Get 函数类型实现某一个接口，称之为接口型函数，方便使用者在调用时既能够传入函数作为参数，也能够传入实现了该接口的结构体作为参数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 缓存的命名空间
type Group struct {
	name      string              // 名称
	getter    Getter              // 未命中时获取源数据的回调
	mainCache cache               // 并发缓存结构
	peers     PeerPicker          // 具备选择peer节点的能力
	loader    *singleflight.Group // 防止缓存击穿
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()

	group := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

// RegisterPeers 注册peer节点
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// Get 缓存获取（本地 -> 远程 -> 回调函数）
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 优先从本地的缓存中获取
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	// 并发场景下，针对相同的key，load过程只会调用一次
	view, err := g.loader.Do(key, func() (any, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})
	if err == nil {
		return view.(ByteView), nil
	}
	return
}

// 从远程节点获取缓存
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &geecachepb.Request{Group: g.name, Key: key}
	resp := &geecachepb.Response{}
	if err := peer.Get(req, resp); err != nil {
		return ByteView{}, err
	}
	return ByteView{b: resp.Value}, nil
}

// 通过本地回调函数获取缓存
func (g *Group) getLocally(key string) (ByteView, error) {
	// 触发回调函数获取源数据
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

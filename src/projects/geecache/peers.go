package geecache

import (
	"learn-go/src/projects/geecache/geecachepb"
)

// PeerPicker peer节点选择器
type PeerPicker interface {
	// PickPeer 选择peer节点
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter peer节点的缓存获取器
type PeerGetter interface {
	// Get 从peer节点中获取缓存
	Get(in *geecachepb.Request, out *geecachepb.Response) error
}

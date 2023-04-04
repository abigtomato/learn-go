package xclient

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

// SelectMode 负载均衡策略
type SelectMode int

const (
	RandomSelect     SelectMode = iota // 随机策略
	RoundRobinSelect                   // 轮询策略
)

// Discovery 服务发现
type Discovery interface {
	Refresh() error                      // 从注册中心更新服务列表
	Update(servers []string) error       // 手动更新服务列表
	Get(mode SelectMode) (string, error) // 更加负载均衡策略，选择一个服务实例
	GetAll() ([]string, error)           // 返回所有的服务实例
}

// MultiServerDiscovery 手动维护服务实例的服务发现
type MultiServerDiscovery struct {
	r       *rand.Rand   // 随机数生成器
	mu      sync.RWMutex // 读写锁
	servers []string     // 服务列表
	index   int          // 轮询算法的索引值
}

func NewMultiServerDiscovery(servers []string) *MultiServerDiscovery {
	d := &MultiServerDiscovery{
		servers: servers,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	d.index = d.r.Intn(math.MaxInt32 - 1)
	return d
}

var _ Discovery = (*MultiServerDiscovery)(nil)

func (d *MultiServerDiscovery) Refresh() error {
	return nil
}

func (d *MultiServerDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	return nil
}

func (d *MultiServerDiscovery) Get(mode SelectMode) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	n := len(d.servers)
	if n == 0 {
		return "", errors.New("rpc discovery: no available servers")
	}
	switch mode {
	case RandomSelect:
		return d.servers[d.r.Intn(n)], nil
	case RoundRobinSelect:
		s := d.servers[d.index%n]
		d.index = (d.index + 1) % n
		return s, nil
	default:
		return "", errors.New("rpc discovery: not supported select mode")
	}
}

func (d *MultiServerDiscovery) GetAll() ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	servers := make([]string, len(d.servers), len(d.servers))
	copy(servers, d.servers)
	return servers, nil
}

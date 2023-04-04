package balancer

import (
	"errors"
	"sync"
)

const (
	IPHashBalancer         = "ip-hash"
	ConsistentHashBalancer = "consistent-hash"
	P2CBalancer            = "p2c"
	RandomBalancer         = "random"
	R2Balancer             = "round-robin"
	LeastLoadBalancer      = "least-load"
	BoundedBalancer        = "bounded"
)

var (
	NoHostError                = errors.New("no host")
	AlgorithmNotSupportedError = errors.New("algorithm not supported")
)

type Balancer interface {
	Add(string)                     // 添加主机
	Remove(string)                  // 删除主机
	Balance(string) (string, error) // 负载选择
	Inc(string)                     // 连接数递增
	Done(string)                    // 连接数递减
}

// Factory 工厂方法设计模式
type Factory func([]string) Balancer

var factories = make(map[string]Factory)

func Build(algorithm string, hosts []string) (Balancer, error) {
	factory, ok := factories[algorithm]
	if !ok {
		return nil, AlgorithmNotSupportedError
	}
	return factory(hosts), nil
}

// TemplateBalancer 模板方法设计模式
type TemplateBalancer struct {
	sync.RWMutex
	hosts []string
}

func (r *TemplateBalancer) Add(host string) {
	r.Lock()
	defer r.Unlock()
	for _, h := range r.hosts {
		if h == host {
			return
		}
	}
	r.hosts = append(r.hosts, host)
}

func (r *TemplateBalancer) Remove(host string) {
	r.Lock()
	defer r.Unlock()
	for i, h := range r.hosts {
		if h == host {
			r.hosts = append(r.hosts[:i], r.hosts[i+1:]...)
			return
		}
	}
}

func (r *TemplateBalancer) Balance(_ string) (string, error) {
	panic("implement me")
}

func (r *TemplateBalancer) Inc(_ string) {
	panic("implement me")
}

func (r *TemplateBalancer) Done(_ string) {
	panic("implement me")
}

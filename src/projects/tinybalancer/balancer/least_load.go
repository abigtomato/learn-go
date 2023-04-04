package balancer

import (
	fibHeap "github.com/starwander/GoFibonacciHeap"
)

func init() {
	factories[LeastLoadBalancer] = NewLeastLoad
}

func (h *host) Tag() any {
	return h.name
}

func (h *host) Key() float64 {
	return float64(h.load)
}

// LeastLoad 对于最小负载算法而言，如果把所有主机的负载值动态存入动态数组中，寻找负载最小节点的时间复杂度为O(N)，如果把主机的负载值维护成一个红黑树，那么寻找负载最小节点的时间复杂度为O(logN)，我们这里利用的数据结构叫做 斐波那契堆 ，寻找负载最小节点的时间复杂度为O(1)
type LeastLoad struct {
	TemplateBalancer
	heap *fibHeap.FibHeap
}

func NewLeastLoad(hosts []string) Balancer {
	ll := &LeastLoad{
		TemplateBalancer: TemplateBalancer{hosts: hosts},
		heap:             fibHeap.NewFibHeap(),
	}

	for _, h := range hosts {
		ll.Add(h)
	}

	return ll
}

func (l *LeastLoad) Add(hostName string) {
	l.Lock()
	defer l.Unlock()

	if l.heap.GetValue(hostName) != nil {
		return
	}

	_ = l.heap.InsertValue(&host{hostName, 0})
}

func (l *LeastLoad) Remove(hostName string) {
	l.Lock()
	defer l.Unlock()

	if l.heap.GetValue(hostName) == nil {
		return
	}

	_ = l.heap.Delete(hostName)
}

func (l *LeastLoad) Balance(_ string) (string, error) {
	l.Lock()
	defer l.Unlock()

	if l.heap.Num() == 0 {
		return "", NoHostError
	}

	return l.heap.MinimumValue().Tag().(string), nil
}

func (l *LeastLoad) Inc(hostName string) {
	l.Lock()
	defer l.Unlock()

	if l.heap.GetValue(hostName) == nil {
		return
	}

	h := l.heap.GetValue(hostName)
	h.(*host).load++
	_ = l.heap.IncreaseKeyValue(h)
}

func (l *LeastLoad) Done(hostName string) {
	l.Lock()
	defer l.Unlock()

	if l.heap.GetValue(hostName) == nil {
		return
	}

	h := l.heap.GetValue(hostName)
	h.(*host).load--
	_ = l.heap.DecreaseKeyValue(h)
}

package balancer

import (
	"hash/crc32"
	"math/rand"
	"time"
)

const Salt = "@#!"

func init() {
	factories[P2CBalancer] = NewP2C
}

type host struct {
	name string
	load uint64 // 主机负载量
}

type P2C struct {
	TemplateBalancer
	hosts   []*host
	rnd     *rand.Rand
	loadMap map[string]*host
}

func NewP2C(hosts []string) Balancer {
	p := &P2C{
		TemplateBalancer: TemplateBalancer{hosts: hosts},
		hosts:            []*host{},
		rnd:              rand.New(rand.NewSource(time.Now().UnixNano())),
		loadMap:          make(map[string]*host),
	}
	for _, h := range hosts {
		p.Add(h)
	}
	return p
}

func (r *P2C) Add(hostName string) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.loadMap[hostName]; ok {
		return
	}
	h := &host{
		name: hostName,
		load: 0,
	}
	r.hosts = append(r.hosts, h)
	r.loadMap[hostName] = h
}

func (r *P2C) Remove(host string) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.loadMap[host]; !ok {
		return
	}

	delete(r.loadMap, host)

	for i, h := range r.hosts {
		if h.name == host {
			r.hosts = append(r.hosts[:i], r.hosts[i+1:]...)
			return
		}
	}
}

func (r *P2C) Balance(key string) (string, error) {
	r.Lock()
	defer r.Unlock()

	if len(r.hosts) == 0 {
		return "", NoHostError
	}

	n1, n2 := r.hash(key)
	host := n2
	if r.loadMap[n1].load <= r.loadMap[n2].load {
		host = n1
	}
	return host, nil
}

// 若请求IP为空，P2C均衡器将随机选择两个代理主机节点，最后选择其中负载量较小的节点；
// 若请求IP不为空，P2C均衡器通过对IP地址以及对IP地址加盐进行CRC32哈希计算，则会得到两个32bit的值，将其对主机数量进行取模，即CRC32(IP) % len(hosts) 、CRC32(IP + salt) % len(hosts)， 最后选择其中负载量较小的节点；
func (r *P2C) hash(key string) (string, string) {
	var n1, n2 string
	if len(key) > 0 {
		n1 = r.hosts[crc32.ChecksumIEEE([]byte(key))%uint32(len(r.hosts))].name
		n2 = r.hosts[crc32.ChecksumIEEE([]byte(key+Salt))%uint32(len(r.hosts))].name
		return n1, n2
	}
	n1 = r.hosts[r.rnd.Intn(len(r.hosts))].name
	n2 = r.hosts[r.rnd.Intn(len(r.hosts))].name
	return n1, n2
}

// Inc 通过Inc、Done函数对主机的负载量进行加减操作
func (r *P2C) Inc(host string) {
	r.Lock()
	defer r.Unlock()

	h, ok := r.loadMap[host]
	if !ok {
		return
	}
	h.load++
}

func (r *P2C) Done(host string) {
	r.Lock()
	defer r.Unlock()

	h, ok := r.loadMap[host]
	if !ok {
		return
	}

	if h.load > 0 {
		h.load--
	}
}

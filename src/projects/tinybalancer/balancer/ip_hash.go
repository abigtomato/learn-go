package balancer

import "hash/crc32"

func init() {
	factories[IPHashBalancer] = NewIPHash
}

type IPHash struct {
	TemplateBalancer
}

func NewIPHash(hosts []string) Balancer {
	return &IPHash{TemplateBalancer: TemplateBalancer{hosts: hosts}}
}

func (r *IPHash) Balance(key string) (string, error) {
	r.Lock()
	defer r.Unlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}
	value := crc32.ChecksumIEEE([]byte(key)) % uint32(len(r.hosts))
	return r.hosts[value], nil
}

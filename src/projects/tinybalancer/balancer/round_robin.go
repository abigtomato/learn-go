package balancer

type RoundRobin struct {
	TemplateBalancer
	i uint64
}

func init() {
	factories[R2Balancer] = NewRoundRobin
}

func NewRoundRobin(hosts []string) Balancer {
	return &RoundRobin{
		TemplateBalancer: TemplateBalancer{hosts: hosts},
		i:                0,
	}
}

func (r *RoundRobin) Balance(_ string) (string, error) {
	r.Lock()
	defer r.Unlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}
	host := r.hosts[r.i%uint64(len(r.hosts))]
	r.i++
	return host, nil
}

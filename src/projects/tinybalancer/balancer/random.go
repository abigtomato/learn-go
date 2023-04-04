package balancer

import (
	"math/rand"
	"time"
)

func init() {
	factories[RandomBalancer] = NewRandom
}

type Random struct {
	TemplateBalancer
	rnd *rand.Rand
}

func NewRandom(hosts []string) Balancer {
	return &Random{
		TemplateBalancer: TemplateBalancer{hosts: hosts},
		rnd:              rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Random) Balance(_ string) (string, error) {
	r.Lock()
	defer r.Unlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}
	return r.hosts[r.rnd.Intn(len(r.hosts))], nil
}

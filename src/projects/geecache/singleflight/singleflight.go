package singleflight

import "sync"

// 代表正在进行中，或已结束的请求
type call struct {
	wg  sync.WaitGroup // 等待组
	val any            // 请求结果
	err error          // 请求错误
}

// Group 用于管理不同key的请求
type Group struct {
	mu sync.Mutex       // 防止g.m出现并发读写问题
	m  map[string]*call // key对应请求的映射
}

// Do 针对相同的key，无论Do被调用多少次，fn都只能被调用一次
func (g *Group) Do(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()
	if g.m == nil {
		// 延迟初始化
		g.m = make(map[string]*call)
	}

	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		// 如果请求正在进行中，go程等待
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	// 发起请求前加锁
	c.wg.Add(1)

	// 进入g.m中，代表key的请求正在执行
	g.m[key] = c
	g.mu.Unlock()

	// 调用fn发起请求
	c.val, c.err = fn()
	// 请求结束解锁
	c.wg.Done()

	// 更新g.m，清空key对应请求的执行记录
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}

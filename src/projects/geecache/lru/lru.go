package lru

import "container/list"

type Cache struct {
	maxBytes  int64                         // 最大内存大小
	nBytes    int64                         // 占用的内存大小
	ll        *list.List                    // 双向链表
	cache     map[string]*list.Element      // 数据字典 值指向双向链表中的节点
	OnEvicted func(key string, value Value) // 回调函数
}

// 双向链表的数据类型
type entry struct {
	key   string // 链表中冗余存储key，当淘汰首节点时通过key从字典中删除对应的映射
	value Value  // 值是实现了Value接口的任意类型
}

type Value interface {
	Len() int // 返回所占的内存大小
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 获取缓存
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 先从字典中根据键获取值
	if elem, ok := c.cache[key]; ok {
		// 键存在，则移动到链表尾部
		c.ll.MoveToFront(elem)
		kv := elem.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 缓存淘汰
func (c *Cache) RemoveOldest() {
	// 获取队首节点，淘汰
	elem := c.ll.Back()
	if elem != nil {
		c.ll.Remove(elem)
		kv := elem.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 新增/更新缓存
func (c *Cache) Add(key string, value Value) {
	if elem, ok := c.cache[key]; ok {
		// 更新操作
		c.ll.MoveToFront(elem)
		kv := elem.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
	} else {
		// 新增操作
		elem := c.ll.PushFront(&entry{key, value})
		c.cache[key] = elem
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	// 触发内存淘汰操作
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

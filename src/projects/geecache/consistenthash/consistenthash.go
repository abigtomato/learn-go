package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// Map 一致性哈希结构
type Map struct {
	hash     Hash           // 哈希函数
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环
	hashMap  map[int]string // 虚拟节点和真实节点的映射
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

	// 默认的hash函数
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 计算虚拟节点的hash值并添加到环上
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			// 维护虚拟节点和真实节点的映射
			m.hashMap[hash] = key
		}
	}
	// 哈希环排序
	sort.Ints(m.keys)
}

// Get 节点选择
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	// 在哈希环上顺时针匹配到第一个虚拟节点的下标
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// idx==len(m.keys)的情况应选择m.keys[0]，keys是一个环状结构，通过取余的方式来处理
	idx = idx % len(m.keys)

	// 通过虚拟节点获取真实节点
	return m.hashMap[m.keys[idx]]
}

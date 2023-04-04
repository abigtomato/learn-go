package basic

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

// 创建map
func TestMapCreate(t *testing.T) {
	// 方式一: 创建时赋值（key无序且不能重复，重复会覆盖）
	m := map[string]string{
		"name":    "mouse",
		"course":  "golang",
		"site":    "ionic",
		"quality": "notepad",
	}

	// 方式二: 通过make()分配空间创建map
	m2 := make(map[string]int, 10)

	// 方式三: 声明式创建是不会分配内存空间的
	var m3 map[string]int

	fmt.Printf("m=%v, m2=%v, m3=%v\n", m, m2, m3)
	fmt.Printf("数据: %v, 类型: %T, 地址: %p, 长度: %v\n", m2, m2, m2, len(m2))
}

// map的增删改查
func TestMapCrud(t *testing.T) {
	m := make(map[string]string)

	// 新增key-value
	m["age"] = "18"
	fmt.Printf("m=%v\n", m)

	// 新增的key相同，则为更新
	m["age"] = "20"
	fmt.Printf("m=%v\n", m)

	// 按key删除
	delete(m, "name")
	fmt.Printf("m=%v\n", m)

	// 按key取value，第二个返回值表示是否存在value
	if site, ok := m["quality"]; !ok {
		fmt.Println("key does not exist ......")
	} else {
		fmt.Printf("site=%v\n", site)
	}

	// 清空key
	m = make(map[string]string, 10)
	fmt.Printf("清空后 m=%v\n", m)
}

// 遍历map
func TestMapIterate(t *testing.T) {
	m := make(map[string]string)
	for k, v := range m {
		fmt.Printf("key=%v, value=%v\n", k, v)
	}
}

// map排序
func TestMapSort(t *testing.T) {
	m := map[int]int{
		10: 100, 1: 13, 4: 56, 8: 90,
	}

	// 将全部key取出存入切片
	var keys []int
	for k, _ := range m {
		keys = append(keys, k)
	}

	// 排序切片
	sort.Ints(keys)

	// 按排序后的顺序依次取出value
	for _, k := range keys {
		fmt.Printf("key=%v, value=%v\n", k, m[k])
	}
}

// 词频统计示例
func TestWordCount(t *testing.T) {
	str := "I love my work and I love my family too"

	slice := strings.Fields(str)
	result := make(map[string]int)

	for _, val := range slice {
		if num, ok := result[val]; !ok {
			result[val] = 1
		} else {
			num++
			result[val] = num
		}
	}

	fmt.Println(result)
}

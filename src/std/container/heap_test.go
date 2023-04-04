package container

import (
	"container/heap"
	"fmt"
	"testing"
)

// 数值堆 可用于实现优先队列
type IntHeap []int

// 长度
func (h *IntHeap) Len() int {
	return len(*h)
}

// 比较
func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

// 交换
func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

// 入堆
func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}

// 出堆
func (h *IntHeap) Pop() any {
	elem := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return elem
}

func TestHeap(t *testing.T) {
	var intHeap IntHeap
	heap.Init(&intHeap)
	heap.Push(&intHeap, 5)
	heap.Push(&intHeap, 3)
	heap.Push(&intHeap, 4)
	heap.Push(&intHeap, 2)
	heap.Push(&intHeap, 1)
	fmt.Println(heap.Pop(&intHeap)) // 1
	fmt.Println(heap.Pop(&intHeap)) // 2
	fmt.Println(heap.Pop(&intHeap)) // 3
	fmt.Println(heap.Pop(&intHeap)) // 4
	fmt.Println(heap.Pop(&intHeap)) // 5
}

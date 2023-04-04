package algorithm

import (
	"fmt"
	"testing"
)

// 初始化堆
func initHeap(arr []int) {
	// 将切片转成二叉树模型  实现大根堆
	length := len(arr)
	for i := length/2 - 1; i >= 0; i-- {
		sortHeap(arr, i, length-1)
	}

	// 根节点存储最大值
	for i := length - 1; i > 0; i-- {
		// 如果只剩下根节点和跟节点下的左子节点
		if i == 1 && arr[0] <= arr[i] {
			break
		}

		// 将根节点和叶子节点数据交换
		arr[0], arr[i] = arr[i], arr[0]

		sortHeap(arr, 0, i-1)
	}
}

// 获取堆中最大值放在根节点
func sortHeap(arr []int, startNode int, maxNode int) {
	// 最大值放在根节点
	var max int

	// 定义做左子节点和右子节点
	lChild := startNode*2 + 1
	rChild := lChild + 1

	// 子节点超过比较范围 跳出递归
	if lChild >= maxNode {
		return
	}

	// 左右比较  找到最大值
	if rChild <= maxNode && arr[rChild] > arr[lChild] {
		max = rChild
	} else {
		max = lChild
	}

	// 和跟节点比较
	if arr[max] <= arr[startNode] {
		return
	}

	// 交换数据
	arr[startNode], arr[max] = arr[max], arr[startNode]

	// 递归进行下次比较
	sortHeap(arr, max, maxNode)
}

// 堆排序
func TestHeapSort(t *testing.T) {
	data := []int{2, 1, 6, 8, 3, 5, 9, 4, 7}
	initHeap(data)
	fmt.Println(data)
}

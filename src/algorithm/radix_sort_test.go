package algorithm

import (
	"fmt"
	"testing"
)

func radixSort(arr []int) {
	max := getMax(arr)
	for bit := 1; max/bit > 0; bit *= 10 {
		bitSort(arr, bit)
	}
}

func bitSort(arr []int, bit int) {
	n := len(arr)
	bitCounts := make([]int, 10)

	for i := 0; i < n; i++ {
		num := (arr[i] / bit) % 10
		bitCounts[num]++
	}

	for i := 1; i < 10; i++ {
		bitCounts[i] += bitCounts[i-1]
	}

	tmp := make([]int, 10)

	for i := n - 1; i >= 0; i-- {
		num := (arr[i] / bit) % 10
		tmp[bitCounts[num]-1] = arr[i]
		bitCounts[num]--
	}

	for i := 0; i < n; i++ {
		arr[i] = tmp[i]
	}
}

// 基数排序
func TestRadixSort(t *testing.T) {
	data := []int{4, 2, 8, 0, 5, 7, 1, 3, 9}
	radixSort(data)
	fmt.Println(data)
}

package algorithm

import (
	"fmt"
	"testing"
)

func countingSort(arr []int) []int {
	bucketLen := getMax(arr) + 1
	bucket := make([]int, bucketLen)

	sortedIndex := 0
	length := len(arr)

	for i := 0; i < length; i++ {
		bucket[arr[i]] += 1
	}

	for j := 0; j < bucketLen; j++ {
		for bucket[j] > 0 {
			arr[sortedIndex] = j
			sortedIndex += 1
			bucket[j] -= 1
		}
	}

	return arr
}

func getMax(arr []int) (max int) {
	max = arr[0]
	for _, v := range arr {
		if max < v {
			max = v
		}
	}
	return
}

// 计数排序
func TestCountingSort(t *testing.T) {
	data := []int{4, 2, 8, 0, 5, 7, 1, 3, 9}
	result := countingSort(data)
	fmt.Println(result)
}

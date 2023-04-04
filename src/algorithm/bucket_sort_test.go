package algorithm

import (
	"fmt"
	"testing"
)

func bucketSort(arr []int) {
	num := len(arr)
	max := getMax(arr)
	buckets := make([][]int, num)

	// 分配入桶
	index := 0
	for i := 0; i < num; i++ {
		index = arr[i] * (num - 1) / max
		buckets[index] = append(buckets[index], arr[i])
	}

	// 桶内排序
	tmpPos := 0
	for i := 0; i < num; i++ {
		bucketLen := len(buckets[i])
		if bucketLen > 0 {
			sortInBucket(buckets[i])
			copy(arr[tmpPos:], buckets[i])
			tmpPos += bucketLen
		}
	}
}

// 此处的桶内排序通过插入排序实现（可以用任意其他排序方式）
func sortInBucket(bucket []int) {
	length := len(bucket)
	if length == 1 {
		return
	}

	for i := 1; i < length; i++ {
		backup := bucket[i]
		j := i - 1
		// 将选出的被排数比较后插入左边有序区
		// 注意j >= 0必须在前边，否则会数组越界
		for j >= 0 && backup < bucket[j] {
			// 移动有序数组
			bucket[j+1] = bucket[j]
			// 反向移动下标
			j--
		}
		// 插队插入移动后的空位
		bucket[j+1] = backup
	}
}

// 桶排序
func TestBucketSort(t *testing.T) {
	data := []int{4, 2, 8, 0, 5, 7, 1, 3, 9}
	bucketSort(data)
	fmt.Println(data)
}

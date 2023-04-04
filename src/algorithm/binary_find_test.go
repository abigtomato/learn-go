package algorithm

import (
	"fmt"
	"testing"
)

func binaryFind(data []int, left int, right int, value int) {
	// 左下标超过右下标时递归结束，表示无法找到
	if left > right {
		return
	}

	// 中间下标，将查找区间分为前后两部分
	middle := (left + right) / 2

	if value > data[middle] {
		// 若是查找的数大于中间的数据，那么将左下标移动到中间下标的后一位，缩短最大查找范围为后半部分
		binaryFind(data, middle+1, right, value)
	} else if value < data[middle] {
		// 若是查找的数小于中间的数据，那么将右下标移动到中间下标的前一位，缩短最大查找范围为前半部分
		binaryFind(data, left, middle-1, value)
	} else {
		fmt.Println(middle)
	}
}

// 二分查找
func TestBinaryFind(t *testing.T) {
	data := []int{1, 8, 10, 89, 1000, 1234}
	binaryFind(data, 0, len(data)-1, 1000)
}

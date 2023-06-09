package algorithm

import (
	"fmt"
	"testing"
)

// 选择排序
//
// 初始状态:
//
//	[8, 3, 2, 1, 7, 4, 6, 5]
//
// 第1次比较并交换:
//
//	[1, 3, 2, 8, 7, 4, 6, 5]
//
// 第2次比较并交换:
//
//	[1, 2, 3, 8, 7, 4, 6, 5]
//
// 第3次比较并交换:
//
//	[1, 2, 3, 4, 7, 8, 6, 5]
//
// 第4次比较并交换:
//
//	[1, 2, 3, 4, 5, 8, 6, 7]
//
// 第5次比较并交换:
//
//	[1, 2, 3, 4, 5, 6, 8, 7]
//
// 第6次比较并交换:
//
//	[1, 2, 3, 4, 5, 6, 7, 8]
//
// 第7次比较并交换:
//
//	[1, 2, 3, 4, 5, 6, 7, 8]
//
// 分析 (升序):
// 1. 第1次从arr[0] ~ arr[n-1]中选取最小值，与arr[0]交换
// 2. 第2次从arr[1] ~ arr[n-1]中选取最小值，与arr[1]交换
// 3. 第3次从arr[2] ~ arr[n-1]中选取最小值，与arr[2]交换
// 4. 第i次从arr[i-1] ~ arr[n-1]中选取最小值，与arr[i-1]交换
// 5. 第n-1次从arr[n-2] ~ arr[n-1]中选取最小值，与arr[n-2]交换
// 6. 总共通过n-1次，得到一个按排序码从小到大排序的有序序列
func selectSort(data []int) {
	// 1.外层遍历整个序列
	for i := 0; i < len(data); i++ {
		// 2.默认0号元素为最大值
		max := 0

		// 3.遍历剩下的元素，找出最大值
		for j := 1; j < len(data)-i; j++ {
			if data[j] > data[max] {
				max = j
			}
		}

		// 4.将最大值移动到序列后面（序列末尾确定一个元素位置后，下次遍历就可以忽略）
		data[max], data[len(data)-1-i] = data[len(data)-1-i], data[max]
	}
}

func TestSelectSort(t *testing.T) {
	data := []int{4, 2, 8, 0, 5, 7, 1, 3, 9}
	selectSort(data)
	fmt.Println(data)
}

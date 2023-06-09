package algorithm

import (
	"fmt"
	"testing"
)

// 插入排序
//
// 初始状态:
//
//	原始表: [23, 0, 12, 56, 34]
//	有序表: [23]
//	无序表: {0, 12, 56, 34}
//
// 第一次找到插入位置:
//
//	[23, 0] {12, 56, 34}
//
// 第二次找到插入位置:
//
//	[23, 0, 0] {56, 34}	insertVal = 12
//	[23, 12, 0]	{56, 34}
//
// 第三次找到插入位置:
//
//	[23, 12, 0, 0] {34} insertVal = 56
//	[23, 12, 12, 0] {34} insertVal = 56
//	[23, 23, 12, 0] {34} insertVal = 56
//	[56, 23, 12, 0] {34}
//
// 第四次找到插入位置:
//
//	[56, 23, 12, 0, 0] insertVal = 34
//	[56, 23, 12, 12, 0] insertVal = 34
//	[56, 23, 23, 12, 0] insertVal = 34
//	[56, 34, 23, 12, 0]
//
// 分析:
// 1. 把n个待排序的元素看成一个有序表和一个无序表
// 2. 有序表开始时只有一个元素，默认是取n个元素的第一个(下标为0)
// 3. 无序表开始时有n-1个元素，除第一个元素外的所有元素
// 4. 排序过程中每次从无序表中取出第一个元素，把它与有序表元素进行比较，将它插入有序表的适当位置，成为新的有序表
func insertSort(data []int) {
	// 外层for遍历无序表，i表示无序表第一个元素（初始时：0为有序表，1~len-1为无序表）
	for i := 1; i < len(data); i++ {
		// j表示无序表第一个元素的前一个元素，也就是有序表最后一个元素（有序表无序表连在一起）
		j := i - 1
		// 每次循环从无序表取出的第一个元素，也就是本次循环需要插入有序表的新元素
		temp := data[i]

		// 新元素从有序表尾部开始循环比较，直到找到适合的插入位置
		for j >= 0 && data[j] > temp {
			// 为新元素的插入腾出位置
			data[j+1] = data[j]
			// 无序表中向前移动
			j--
		}

		// 插入到有序表
		data[j+1] = temp
	}
}

// 插入排序
func TestInsertSort(t *testing.T) {
	data := []int{4, 2, 8, 0, 5, 7, 1, 3, 9}
	insertSort(data)
	fmt.Println(data)
}

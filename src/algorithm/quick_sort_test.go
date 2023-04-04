package algorithm

import (
	"fmt"
	"testing"
)

// 快速排序
//
// 初始状态:
//
//	l = 0, r = 5, pivot = 2
//	[-9, 78, 0, 23, -567, 7]
//
// 第一次排序:
//
//	l = 1, r = 4, arr[l] = 78, arr[r] = -567
//	[-9, -567, 0, 23, 78, 7]
//	当l = 3, r = 1, arr[l] = 23, arr[r] = -567时，l >= r 退出第一次排序
//
// 左递归:
//
//	QuickSort(left, r, arr)
//	left = 0, r = 1, arr = [-9, -567]
//
// 右递归:
//
//	QuickSort(l, right, arr)
//	l = 3, right = 5, arr = [23, 78, 7]
//
// 分析:
// 1. 通过第一次排序将要排序的数据按照中间值分割成独立的两部分
// 2. 期望的情况是左边部分的所有数据都比中间值要小，右边部分要比中间大
// 3. 之后再次按照此方法对这两部分数据分别进行快速排序，整个排序过程可以递归进行
func quickSort(left, right int, data []int) {
	// 左右指针
	l, r := left, right
	// 基准值
	pivot := data[(left+right)/2]

	// 1.for的作用就是把比pivot小的数移到左边，比pivot大的数移到右边
	for l < r {
		// 2.从左边找一个比pivot大的值
		for data[l] < pivot {
			l++
		}

		// 3.从右边找一个比pivot小的值
		for data[r] > pivot {
			r--
		}

		// 4.交换数据
		data[l], data[r] = data[r], data[l]
	}

	// 5.左右指针指向同一个数据则各走一步分开
	if l == r {
		l++
		r--
	}

	// 6.向左边递归
	if left < r {
		quickSort(left, r, data)
	}

	// 7.向右边递归
	if right > l {
		quickSort(l, right, data)
	}
}

func TestQuickSort(t *testing.T) {
	data := []int{4, 2, 8, 0, 5, 7, 1, 3, 9}
	quickSort(0, len(data)-1, data)
	fmt.Println(data)
}

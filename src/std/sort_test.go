package std

import (
	"fmt"
	"sort"
	"testing"
)

type TestSlice [][]int

func (t TestSlice) Len() int { return len(t) }

func (t TestSlice) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t TestSlice) Less(i, j int) bool { return t[i][1] < t[j][1] }

type TestSlice2 []map[string]float64

func (t TestSlice2) Len() int { return len(t) }

func (t TestSlice2) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t TestSlice2) Less(i, j int) bool { return t[i]["a"] < t[j]["a"] }

type People struct {
	Name string
	Age  int
}

type TestSlice3 []People

func (t TestSlice3) Len() int { return len(t) }

func (t TestSlice3) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t TestSlice3) Less(i, j int) bool { return t[i].Age < t[j].Age }

func TestSort(t *testing.T) {
	s := []int{2, 4, 1, 3}
	sort.Ints(s)
	fmt.Printf("s: %v\n", s)

	f := []float64{1.1, 4.4, 5.5, 3.3, 2.2}
	sort.Float64s(f)
	fmt.Printf("f: %v\n", f)

	// 数字字符串 根据第一个字符大小排序
	// StringSlice等价于[]string切片
	ls := sort.StringSlice{"100", "42", "41", "3", "2"}
	sort.Strings(ls)
	fmt.Printf("ls: %v\n", ls)

	// 字母字符串 根据第一个字符大小排序
	ls = sort.StringSlice{"d", "ac", "c", "ab", "e"}
	sort.Strings(ls)
	fmt.Printf("ls: %v\n", ls)

	// 汉字字符串 比较byte大小
	ls = sort.StringSlice{"啊", "博", "次", "得", "饿", "周"}
	sort.Strings(ls)
	fmt.Printf("ls: %v\n", ls)

	// 按第二个位置排序
	ts := TestSlice{{1, 4}, {9, 3}, {7, 5}}
	sort.Sort(ts)
	fmt.Printf("ts: %v\n", ts)

	// 按照"a"对应的值排序
	ts2 := TestSlice2{{"a": 4, "b": 12}, {"a": 3, "b": 11}, {"a": 5, "b": 10}}
	sort.Sort(ts2)
	fmt.Printf("ts2: %v\n", ts2)

	// 按照Age字段排序
	ts3 := TestSlice3{{Name: "n1", Age: 12}, {Name: "n2", Age: 11}, {Name: "n3", Age: 10}}
	sort.Sort(ts3)
	fmt.Printf("ts3: %v\n", ts3)
}

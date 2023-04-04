package basic

import (
	"fmt"
	"testing"
)

// ID和类型int64是不同的
type ID int64

// 声明了一个叫做NumberDerived的类型约束
type NumberDerived interface {
	// 联合了int64和float64的类型约束，使用 | 去表示联合
	// 使用~int64表示int64和它的任何衍生类型
	~int64 | ~float64
}

// comparable已经被Go预先声明（在builtin.go文件里可以找到），它表示任何能够使用==和!=进行比较的类型
func SumNumbersDerived[K comparable, V NumberDerived](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func ForEach[T any](list []T, action func(T)) {
	for _, item := range list {
		action(item)
	}
}

// 泛型
func TestGeneric(t *testing.T) {
	intMap := map[string]int64{
		"first":  34,
		"second": 12,
	}
	floatMap := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}
	idMap := map[string]ID{
		"first":  ID(34),
		"second": ID(12),
	}

	fmt.Printf("SumNumbersDerived(intMap): %v\n", SumNumbersDerived(intMap))
	fmt.Printf("SumNumbersDerived(floatMap): %v\n", SumNumbersDerived(floatMap))
	fmt.Printf("SumNumbersDerived(idMap): %v\n", SumNumbersDerived(idMap))

	ForEach([]string{"hello", "golang"}, func(s string) {
		fmt.Printf("%v\n", s)
	})
}

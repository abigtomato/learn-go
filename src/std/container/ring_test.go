package container

import (
	"container/ring"
	"fmt"
	"testing"
)

func TestRing(t *testing.T) {
	// 只指定我们需要存储的值，不需要管前后指针的初始化，因为Ring的方法都会在执行前检查Ring是否初始化，如果没有初始化则先进行初始化
	r1 := ring.Ring{Value: 1}

	// New()的参数是指定要创建圈元素的个数，而且使用New()前后指针是已经初始化好的
	r2 := ring.New(1)
	r2.Value = 2

	// 链接两个单圈
	r1.Link(r2)

	// unlink掉后面1个节点
	r1.Unlink(1)

	// 分别是前一个节点、后一个节点、移动多少个节点（Move(1)==Next()，Move(-1)==Prev())
	cur := &r1
	fmt.Println(cur.Value)
	cur = cur.Prev()
	fmt.Println(cur.Value)
	cur = cur.Next()
	fmt.Println(cur.Value)
	cur = cur.Move(1)
	fmt.Println(cur.Value)

	// 遍历整个圈
	r1.Do(func(a any) {
		fmt.Println(a)
	})
}

package datastructure

import (
	"fmt"
	"reflect"
	"testing"
)

type TreeNode struct {
	Data  any
	Left  *TreeNode
	Right *TreeNode
}

// 前序遍历
func (n *TreeNode) preOrder() {
	if n == nil {
		return
	}
	fmt.Println(n.Data)
	n.Left.preOrder()
	n.Right.preOrder()
}

// 中序遍历
func (n *TreeNode) midOrder() {
	if n == nil {
		return
	}
	n.Left.midOrder()
	fmt.Println(n.Data)
	n.Right.midOrder()
}

// 后序遍历
func (n *TreeNode) rearOrder() {
	if n == nil {
		return
	}
	n.Left.rearOrder()
	n.Right.rearOrder()
	fmt.Println(n.Data)
}

// 树的高度
func (n *TreeNode) height() int {
	if n == nil {
		return 0
	}

	lh := n.Left.height()
	rh := n.Right.height()

	if lh > rh {
		lh++
		return lh
	} else {
		rh++
		return rh
	}
}

// 叶子节点树
func (n *TreeNode) leafCount(num *int) {
	if n == nil {
		return
	}

	if n.Left == nil && n.Right == nil {
		*num++
	}

	n.Left.leafCount(num)
	n.Right.leafCount(num)
}

// 二分搜索
func (n *TreeNode) search(data any) {
	if n == nil {
		return
	}

	if reflect.TypeOf(n.Data) == reflect.TypeOf(data) && n.Data == data {
		fmt.Println("数据存在: ", data)
		return
	}

	n.Left.search(data)
	n.Right.search(data)
}

// 销毁树
func (n *TreeNode) destroy() {
	if n == nil {
		return
	}

	n.Left.destroy()
	n.Left = nil

	n.Right.destroy()
	n.Right = nil

	n.Data = nil
}

// 反转树
func (n *TreeNode) reverse() {
	if n == nil {
		return
	}

	n.Left, n.Right = n.Right, n.Left

	n.Left.reverse()
	n.Right.reverse()
}

// 拷贝树
func (n *TreeNode) copy() *TreeNode {
	if n == nil {
		return nil
	}

	left := n.Left.copy()
	right := n.Right.copy()

	return &TreeNode{
		Data:  n.Data,
		Left:  left,
		Right: right,
	}
}

// 二叉树
func TestBinaryTree(t *testing.T) {
	node := &TreeNode{
		Data: 0,
		Left: &TreeNode{
			Data:  1,
			Left:  &TreeNode{Data: 3},
			Right: &TreeNode{Data: 4},
		},
		Right: &TreeNode{
			Data:  2,
			Left:  &TreeNode{Data: 5},
			Right: &TreeNode{Data: 6},
		},
	}

	node.preOrder()
	node.midOrder()
	node.rearOrder()

	height := node.height()
	fmt.Println(height)

	num := 0
	node.leafCount(&num)
	fmt.Println(num)

	node.search(1)

	node.reverse()
	node.preOrder()

	newNode := node.copy()
	newNode.preOrder()
}

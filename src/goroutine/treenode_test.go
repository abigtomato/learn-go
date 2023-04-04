package goroutine

import (
	"fmt"
	"testing"
)

// 树节点
type TreeNode struct {
	Value       int
	Left, Right *TreeNode
}

func (n *TreeNode) Print() {
	fmt.Printf("[%v] ", n.Value)
}

// 函数式编程（接收函数参数）
func (n *TreeNode) TraverseFunc(fun func(*TreeNode)) {
	if n == nil {
		return
	}

	n.Left.TraverseFunc(fun)
	fun(n)
	n.Right.TraverseFunc(fun)
}

// 转换为节点channel
func (n *TreeNode) TraverseWithChannel() chan *TreeNode {
	nodeChan := make(chan *TreeNode)

	// 开启协程遍历二叉树节点并入队channel
	go func(nodeChan chan *TreeNode) {
		defer close(nodeChan)
		n.TraverseFunc(func(node *TreeNode) {
			nodeChan <- node
		})
	}(nodeChan)

	return nodeChan
}

// 使用goroutine + channel遍历二叉树
func TestTreeNode(t *testing.T) {
	// 形成一颗二叉树
	root := &TreeNode{
		Value: 3,
		Left: &TreeNode{
			Value: 0,
			Right: &TreeNode{
				Value: 2,
			},
		},
		Right: &TreeNode{
			Value: 5,
			Left: &TreeNode{
				Value: 4,
			},
		},
	}

	// 将tree转换为node channel
	nodeChan := root.TraverseWithChannel()

	// 遍历node获取最大值的节点
	maxNode := &TreeNode{}
	for node := range nodeChan {
		fmt.Printf("Node Value: %v\n", node.Value)
		if node.Value > maxNode.Value {
			maxNode = node
		}
	}
	fmt.Printf("Max Node Value: %v\n", maxNode.Value)
}

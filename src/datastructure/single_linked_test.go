package datastructure

import (
	"fmt"
	"testing"
)

type SingleNode struct {
	no   int
	name string
	next *SingleNode // 后继指针
}

type SingleLinkedList struct {
	// 头节点
	head *SingleNode
}

// 尾部追加节点
func (s *SingleLinkedList) appendNode(newHeroNode *SingleNode) {
	// 临时指针，用于后移遍历整个链表的节点
	temp := s.head

	for {
		// 若当前指向节点的后继指针指向nil，则代表链表已遍历至末尾
		if temp.next == nil {
			break
		}

		// 向后移动指针
		temp = temp.next
	}

	// 执行到此则代表已经退出循环，temp指针指向的是最后一个节点，直接将新节点挂在后面即可
	temp.next = newHeroNode
}

// 从链表的中间插入节点（按编号排序）
func (s *SingleLinkedList) insertNode(newHeroNode *SingleNode) {
	temp := s.head

	for {
		if temp.next == nil {
			break
		} else if temp.next.no >= newHeroNode.no {
			// 断开前后两节点的连接，中间插入新节点并重新建立连接
			newHeroNode.next = temp.next
			temp.next = newHeroNode
			return
		}

		temp = temp.next
	}

	// 执行到此则代表已经退出循环，temp指针指向的是最后一个节点，直接将新节点挂在后面即可
	temp.next = newHeroNode
}

// 根据名称搜索节点
func (s *SingleLinkedList) getNodeByName(name string) (node *SingleNode) {
	temp := s.head

	for {
		if temp.next == nil {
			break
		} else if temp.next.name == name {
			node = temp.next
			return
		}

		temp = temp.next
	}

	return
}

// 删除节点
func (s *SingleLinkedList) delNodeByNo(no int) {
	temp := s.head

	for {
		if temp.next == nil {
			break
		} else if temp.next.no == no {
			// 直接将待删除节点的前一个节点的next指针指向待删除节点的后一个节点即可
			temp.next = temp.next.next
			return
		}

		temp = temp.next
	}
}

// 更新节点信息
func (s *SingleLinkedList) updateNode(updateHeroNode *SingleNode) {
	temp := s.head

	for {
		if temp.next == nil {
			break
		} else if temp.next.no == updateHeroNode.no {
			updateHeroNode.next = temp.next.next
			temp.next = updateHeroNode
			return
		}

		temp = temp.next
	}
}

// 展示链表数据
func (s *SingleLinkedList) showLinkedList() {
	temp := s.head

	if temp.next == nil {
		fmt.Printf("IsEmpty(head) fail, linkedlist empty")
		return
	}

	for {
		if temp.next == nil {
			break
		}
		fmt.Printf("[val: %v, ptr: %p]--->", (*temp).next.name, (*temp).next.next)
		temp = temp.next
	}
}

// 单链表
func TestSingleLinkedList(t *testing.T) {
	linkedList := SingleLinkedList{
		head: &SingleNode{},
	}

	linkedList.appendNode(&SingleNode{no: 1, name: "albert"})
	linkedList.appendNode(&SingleNode{no: 2, name: "lily"})
	linkedList.appendNode(&SingleNode{no: 5, name: "charname"})

	linkedList.insertNode(&SingleNode{no: 4, name: "yahaha"})

	linkedList.updateNode(&SingleNode{no: 2, name: "Aliah"})

	linkedList.delNodeByNo(1)

	fmt.Println(linkedList.getNodeByName("charname"))

	linkedList.showLinkedList()
}

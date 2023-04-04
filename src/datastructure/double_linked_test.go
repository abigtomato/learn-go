package datastructure

import (
	"fmt"
	"testing"
)

type DoubleNode struct {
	no   int
	name string
	pre  *DoubleNode // 前趋指针
	next *DoubleNode // 后继指针
}

type DoubleLinkedList struct {
	head *DoubleNode
}

// 链表末尾追加节点
func (d *DoubleLinkedList) appendNode(newNode *DoubleNode) {
	// 临时指针，用于遍历节点
	var temp *DoubleNode
	temp = d.head

	for {
		if (*temp).next == nil {
			break
		}

		temp = (*temp).next
	}

	// 此时temp指向最后一个节点，使temp指向节点后继指向新节点，新节点前序指向temp指向节点
	(*temp).next = newNode
	newNode.pre = temp
}

// 插入节点
func (d *DoubleLinkedList) insertNode(newNode *DoubleNode) {
	var temp *DoubleNode
	temp = d.head

	for {
		if (*temp).next == nil {
			break
		} else if (*temp).next.no >= newNode.no {
			// 断开原节点的连接，插入新节点，重新建立连接
			newNode.pre = (*temp).pre
			(*temp).pre.next = newNode

			(*temp).pre = newNode
			newNode.next = temp

			return
		}

		temp = (*temp).next
	}

	// 节点插入末尾的情况
	(*temp).next = newNode
	newNode.pre = temp
}

// 删除节点
func (d *DoubleLinkedList) delNodeByName(name string) {
	var temp *DoubleNode
	temp = d.head

	for {
		if (*temp).next == nil {
			if (*temp).name == name {
				(*temp).pre.next = nil
			}
			break
		} else if (*temp).name == name {
			(*temp).pre.next = (*temp).next
			(*temp).next.pre = (*temp).next

			return
		}

		temp = (*temp).next
	}
}

// 展示链表所有节点
func (d *DoubleLinkedList) showLinkedList() {
	var temp *DoubleNode
	temp = d.head

	for {
		if (*temp).next == nil {
			break
		}

		fmt.Printf("[no=%v, name=%v, pre=%p, next=%p]-->", (*temp).next.no, (*temp).next.name, (*temp).next.pre, (*temp).next.next)
		temp = (*temp).next
	}
}

// 判空
func (d *DoubleLinkedList) isEmpty() bool {
	return d.head.next == nil
}

// 双向链表
func TestDoubleLinkedList(t *testing.T) {
	link := &DoubleLinkedList{
		head: &DoubleNode{},
	}

	link.appendNode(&DoubleNode{no: 1, name: "hadoop"})
	link.appendNode(&DoubleNode{no: 2, name: "spark"})
	link.appendNode(&DoubleNode{no: 3, name: "kafka"})
	link.appendNode(&DoubleNode{no: 5, name: "storm"})
	link.showLinkedList()
	fmt.Println()

	link.insertNode(&DoubleNode{no: 4, name: "hbase"})
	link.insertNode(&DoubleNode{no: 6, name: "hive"})
	link.showLinkedList()
	fmt.Println()

	link.delNodeByName("hive")
	link.delNodeByName("kafka")
	link.showLinkedList()
	fmt.Println()
}

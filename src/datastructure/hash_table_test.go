package datastructure

import (
	"fmt"
	"testing"
)

type EmpNode struct {
	Id   int
	Name string
	Next *EmpNode
}

type EmpLink struct {
	Head *EmpNode
}

func (l *EmpLink) insert(emp *EmpNode) {
	if l.Head == nil {
		l.Head = emp
		return
	}

	var cur, pre *EmpNode
	cur = l.Head
	pre = nil

	for cur != nil {
		if cur.Id >= emp.Id {
			break
		}
		pre = cur
		cur = cur.Next
	}
	pre.Next = emp
	emp.Next = cur
}

func (l *EmpLink) showList() {
	if l.Head == nil {
		fmt.Println("linked list empty")
		return
	}

	cur := l.Head
	for cur != nil {
		fmt.Printf("[Id: %v, Name: %v, Next: %p]--->", cur.Id, cur.Name, cur.Next)
		cur = cur.Next
	}
}

func (l *EmpLink) findEmpById(id int) *EmpNode {
	if l.Head == nil {
		fmt.Println("linked list empty")
		return nil
	}

	cur := l.Head
	for cur != nil {
		if cur.Id == id {
			return cur
		}
		cur = cur.Next
	}
	return nil
}

type HashTable struct {
	LinkArr [7]EmpLink
}

func (t *HashTable) insert(emp *EmpNode) {
	insertNo := t.hashFun(emp.Id)
	t.LinkArr[insertNo].insert(emp)
}

func (t *HashTable) hashFun(id int) int {
	return id % (len(t.LinkArr) - 1)
}

func (t *HashTable) showList() {
	for index, value := range t.LinkArr {
		fmt.Printf("linked list %v number: ", index)
		value.showList()
		fmt.Println()
	}
}

func (t *HashTable) findEmpById(id int) *EmpNode {
	findNo := t.hashFun(id)
	return t.LinkArr[findNo].findEmpById(id)
}

// 哈希表
func TestHashTable(t *testing.T) {
	hashTable := &HashTable{}

	for i := 6; i <= 54; i += 6 {
		emp := &EmpNode{
			Id:   i,
			Name: fmt.Sprintf("albert %v", i),
		}
		hashTable.insert(emp)
	}

	hashTable.showList()

	res := hashTable.findEmpById(42)
	fmt.Printf("result: %v\n", res.Name)
}

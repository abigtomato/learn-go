package datastructure

import (
	"errors"
	"fmt"
	"testing"
)

type CircleQueue struct {
	MaxSize int
	Values  []interface{}
	Head    int
	Tail    int
}

func NewCircleQueue(maxSize int) *CircleQueue {
	return &CircleQueue{
		MaxSize: maxSize,
		Values:  make([]interface{}, maxSize),
		Head:    0,
		Tail:    0,
	}
}

// 入队
func (c *CircleQueue) push(val interface{}) (err error) {
	if c.isFull() {
		err = errors.New("push fail queue full")
		return
	}

	// 下标每次移动后都要与队列最大长度取模
	// 这样才会达成循环的目的，移动到尾部后和maxsize取模会再次指向首部
	c.Values[c.Tail] = val
	c.Tail = (c.Tail + 1) % c.MaxSize

	return
}

// 出队
func (c *CircleQueue) pop() (val interface{}, err error) {
	if c.isEmpty() {
		err = errors.New("pop fail queue empty")
		return
	}

	val = c.Values[c.Head]
	c.Head = (c.Head + 1) % c.MaxSize

	return
}

// 计算队列元素数量
func (c *CircleQueue) show() {
	size := c.size()

	if size == 0 {
		fmt.Println("queue size = 0")
		return
	}

	// 临时的指向，从头下标开始不断后移遍历所有元素
	tempHead := c.Head
	for i := 0; i < size; i++ {
		fmt.Printf("queue[%v]=%v\n", tempHead, c.Values[tempHead])
		tempHead = (tempHead + 1) % c.MaxSize
	}
}

// 判断队列是否为空
func (c *CircleQueue) isEmpty() bool {
	return c.Head == c.Tail
}

// 判断队列是否已满
func (c *CircleQueue) isFull() bool {
	return (c.Tail+1)%c.MaxSize == c.Head
}

// 队列元素个数
func (c *CircleQueue) size() int {
	return (c.Tail + c.MaxSize - c.Head) % c.MaxSize
}

// 使用数组实现循环队列
func TestCircleQueue(t *testing.T) {
	queue := NewCircleQueue(10)

	for i := 1; i <= 10; i++ {
		_ = queue.push(i)
	}
	queue.show()

	for i := 0; i <= 5; i++ {
		val, err := queue.pop()
		if err != nil {
			break
		}
		fmt.Printf("pop val=%v\n", val)
	}

	for i := 1; i <= 5; i++ {
		_ = queue.push(i)
	}
	queue.show()
}

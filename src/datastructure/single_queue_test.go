package datastructure

import (
	"errors"
	"fmt"
	"testing"
)

type SingleQueue struct {
	MaxSize int           // 最大容量
	Values  []interface{} // 底层存储结构
	Font    int           // 头下标指向
	Rear    int           // 尾下标指向
}

// NewQueue 初始化一个队列
func NewQueue(maxSize int) (queue *SingleQueue) {
	queue = &SingleQueue{
		MaxSize: maxSize,
		Values:  make([]interface{}, maxSize),
		Font:    -1,
		Rear:    -1,
	}
	return
}

// 入队
func (q *SingleQueue) push(value interface{}) (err error) {
	// 若尾下标指向了底层数组最后一个元素，则代表队列已满
	if q.Rear == q.MaxSize-1 {
		err = errors.New("add fail queue full")
		return
	}

	// 尾下标后移，从队尾入队
	q.Rear++
	q.Values[q.Rear] = value

	return
}

// 出队
func (q *SingleQueue) pop() (val interface{}, err error) {
	// 若头下标和尾下标指向相同，则代表队列为空
	if q.Font == q.Rear {
		err = errors.New("get fail queue empty")
		return
	}

	// 头下标后移，从队头出队
	q.Font++
	val = q.Values[q.Font]

	return
}

// 显示队列数据
func (q *SingleQueue) show() {
	if q.Font == q.Rear {
		return
	}

	// 队列的数据范围从头下标到尾下标
	for i := q.Font; i <= q.Rear; i++ {
		fmt.Printf("queue[%v]=%v\n", i, q.Values[i])
	}
}

// 使用数组实现单向队列
func TestSingleQueue(t *testing.T) {
	queue := NewQueue(10)

	for i := 1; i <= 5; i++ {
		err := queue.push(i)
		if err != nil {
			return
		}
	}

	if val, err := queue.pop(); err != nil {
		fmt.Printf("queue.Get() fail error: %v\n", err.Error())
		return
	} else {
		fmt.Printf("val=%v\n", val)
	}

	queue.show()
}

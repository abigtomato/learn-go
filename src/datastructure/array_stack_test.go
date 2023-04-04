package datastructure

import (
	"errors"
	"fmt"
	"testing"
)

type ArrayStack struct {
	MaxSize int
	Top     int
	Values  []any
}

func (s *ArrayStack) push(value any) error {
	if s.Top == s.MaxSize-1 {
		return errors.New("push fail. stack pull")
	}

	s.Top++
	s.Values[s.Top] = value

	return nil
}

func (s *ArrayStack) pop() (value any, err error) {
	if s.Top == -1 {
		err = errors.New("pop fail. stack empty")
		return
	}

	value = s.Values[s.Top]
	s.Top--

	return
}

func (s *ArrayStack) list() error {
	if s.Top == -1 {
		return errors.New("list fail. stack empty")
	}

	for i := s.Top; i >= 0; i-- {
		fmt.Printf("->%v\n", s.Values[i])
	}

	return nil
}

// 数组栈
func TestArrayStack(t *testing.T) {
	stack := &ArrayStack{
		MaxSize: 5,
		Top:     -1,
		Values:  make([]any, 5),
	}

	_ = stack.push(1)
	_ = stack.push(2)
	_ = stack.push(3)
	_ = stack.push(4)
	_ = stack.push(5)

	_ = stack.list()

	val, _ := stack.pop()
	fmt.Println("pop: ", val)
}

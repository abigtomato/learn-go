package datastructure

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

type Stack struct {
	MaxSize int
	Top     int
	Values  []interface{}
}

func (s *Stack) push(value interface{}) error {
	if s.Top == s.MaxSize-1 {
		return errors.New("push fail. stack pull")
	}

	s.Top++
	s.Values[s.Top] = value

	return nil
}

func (s *Stack) pop() (value interface{}, err error) {
	if s.Top == -1 {
		err = errors.New("pop fail. stack empty")
		return
	}

	value = s.Values[s.Top]
	s.Top--

	return
}

func (s *Stack) list() error {
	if s.Top == -1 {
		return errors.New("list fail. stack empty")
	}

	for i := s.Top; i >= 0; i-- {
		fmt.Printf("->%v\n", s.Values[i])
	}

	return nil
}

// 是否为操作符
func isOperate(val int) bool {
	if val == 42 || val == 43 || val == 45 || val == 47 {
		return true
	} else {
		return false
	}
}

// 是否为数字
func isNum(val int) bool {
	if val >= 48 && val <= 57 {
		return true
	} else {
		return false
	}
}

// 计算
func cal(num1 int, num2 int, operate int) (res int, err error) {
	switch operate {
	case 42:
		res = num2 * num1
	case 43:
		res = num2 + num1
	case 45:
		res = num2 - num1
	case 47:
		res = num2 / num1
	default:
		err = errors.New("oper error")
	}
	return
}

// 优先级
func priority(operate int) (res int) {
	if operate == 42 || operate == 47 {
		res = 1
	} else if operate == 43 || operate == 45 {
		res = 0
	}
	return
}

// 提交计算式
func calculation(numStack, operateStack *Stack, exp string) (res int, err error) {
	chExp := []byte(exp)
	var keepNum string

	for i, ch := range chExp {
		if isOperate(int(ch)) {
			if operateStack.Top == -1 {
				err := operateStack.push(int(ch))
				if err != nil {
					return 0, err
				}
			} else {
				if priority(int(ch)) >= priority(operateStack.Values[operateStack.Top].(int)) {
					err := operateStack.push(int(ch))
					if err != nil {
						return 0, err
					}
				} else {
					num1, _ := numStack.pop()
					num2, _ := numStack.pop()
					operate, _ := operateStack.pop()
					res, err = cal(num1.(int), num2.(int), operate.(int))
					err := numStack.push(res)
					if err != nil {
						return 0, err
					}
					err = operateStack.push(int(ch))
					if err != nil {
						return 0, err
					}
				}
			}
		} else if isNum(int(ch)) {
			keepNum += fmt.Sprintf("%c", ch)
			if i == len(chExp)-1 {
				num, _ := strconv.ParseInt(keepNum, 10, 64)
				err := numStack.push(int(num))
				if err != nil {
					return 0, err
				}
			} else {
				if isOperate(int(chExp[i+1])) {
					num, _ := strconv.ParseInt(keepNum, 10, 64)
					err := numStack.push(int(num))
					if err != nil {
						return 0, err
					}
					keepNum = ""
				}
			}
		}
	}

	for operateStack.Top != -1 {
		num1, _ := numStack.pop()
		num2, _ := numStack.pop()
		operate, _ := operateStack.pop()
		res, err = cal(num1.(int), num2.(int), operate.(int))
		err := numStack.push(res)
		if err != nil {
			return 0, err
		}
	}

	num, err := numStack.pop()
	res = num.(int)

	return
}

// 使用栈进行表达式计算
func TestExpStack(t *testing.T) {
	numStack := &Stack{
		MaxSize: 20,
		Top:     -1,
		Values:  make([]interface{}, 20),
	}
	operateStack := &Stack{
		MaxSize: 20,
		Top:     -1,
		Values:  make([]interface{}, 20),
	}

	res, _ := calculation(numStack, operateStack, "30+20*60-200")
	fmt.Printf("result: %v\n", res)
}

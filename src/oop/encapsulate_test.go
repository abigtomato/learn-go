package oop

import (
	"fmt"
	"testing"
)

// 结构体名首字母小写，表示包内私有，只允许包外使用工厂函数创建实例
type Human struct {
	// 封装2: 属性首字母小写，表示为私有，包外不可访问，只允许包外通过对外开放的方法操作
	name string
	age  int
	sal  float64
}

// 提供工厂函数，只能通过工厂函数创建结构体的实例
func NewPerson(name string, age int, sal float64) *Human {
	return &Human{name, age, sal}
}

// 提供包外可见的对结构体属性进行操作的方法
func (p *Human) SetAge(age int) {
	if age >= 150 || age <= 0 {
		return
	}
	p.age = age
}

func (p *Human) GetAge() int {
	return p.age
}

func (p *Human) SetSal(sal float64) {
	if sal > 30000 || sal < 4000 {
		return
	}
	p.sal = sal
}

func (p *Human) GetSal() float64 {
	return p.sal
}

func TestEncapsulate(t *testing.T) {
	person := NewPerson("albert", 20, 4000)
	person.SetSal(5000)
	fmt.Println(person.GetSal())
}

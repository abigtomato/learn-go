package oop

import (
	"fmt"
	"testing"
)

type Person struct {
	Height int
	Weight int
}

type Student struct {
	Name  string
	Age   int
	Score int
}

func (stu *Student) ShowInfo() {
	fmt.Printf("name=%v, age=%v, score=%v\n", stu.Name, stu.Age, stu.Score)
}

func (stu *Student) SetScore(score int) {
	stu.Score = score
}

type Pupil struct {
	// 嵌入了Student匿名结构体
	// 使用匿名结构体实现OOP编程的继承特性，此时Pupil结构体内部存在Student的属性和方法
	// 结构体可以使用嵌套匿名结构体所有的字段和方法(不论大小写)
	Student
	Person
}

type Graduate struct {
	Student
	// 嵌入有名结构体，这种模式称为组合
	Per Person
}

func TestExtends(t *testing.T) {
	// 操作嵌套匿名结构体的属性方法
	pupil := Pupil{}
	pupil.Student.Name = "tom"
	pupil.Student.Age = 8
	pupil.Student.SetScore(30)
	pupil.Student.ShowInfo()

	// 1.编译器会查找结构体中有没有对应的属性和方法，如果没有会继续在嵌套的匿名结构体中查找
	// 2.当结构体和匿名结构体有相同字段和方法时，会将采用就近原则访问，如果希望访问匿名结构体中的数据，则显式调用
	// 3.结构体嵌入多个匿名结构体，如果两个匿名结构体都存在匿名的结构体和方法(同时结构体本身没有)，那么在访问时就必须指定匿名结构体的名字
	graduate := &(Graduate{})
	// 4.如果是组合调用，那么必须通过有名结构体名调用
	graduate.Per.Height = 180
	graduate.Name = "albert"
	graduate.Age = 18
	graduate.SetScore(80)
	graduate.ShowInfo()
	// 5.嵌套匿名结构体后，可以在实例时直接指定匿名结构体字段的值
	graduate2 := Graduate{
		Student: Student{
			Name:  "albert",
			Age:   21,
			Score: 100,
		},
	}
	fmt.Println(graduate2)
}

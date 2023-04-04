package oop

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

type Usb interface {
	// Start 接口中不包含变量，可以定义一组不需要实现的方法体
	// 遵循高内聚低耦合的原则设计
	Start()
	Stop()
}

type Phone struct{}

func (p *Phone) Start() {
	fmt.Println("phone start ...")
}

func (p *Phone) Stop() {
	fmt.Println("phone stop ...")
}

func (p *Phone) Call() {
	fmt.Println("phone call ...")
}

type Camera struct{}

func (c *Camera) Start() {
	fmt.Println("camera start ...")
}

func (c *Camera) Stop() {
	fmt.Println("camera stop ...")
}

// 此处的usb变量是被多种自定义类型实现的接口，具备多种形态，又称多态变量
func Working(usb Usb) {
	usb.Start()
	// usb接口有多种实现，若是其他实现则无法调用phone专属的方法call
	// 直接将phone赋值给usb接口会报错，需要使用类型断言语法转换
	// 此时使用类型断言判断当前传入的usb接口是否指向phone，若是则调用
	if phone, ok := usb.(*Phone); ok {
		phone.Call()
	}
	usb.Stop()
}

// 接口和多态
func TestInterfaceParam(t *testing.T) {
	// 接口不需要显式的实现，只需要变量实现接口定义的所有方法，那么这个变量就实现了这个接口
	phone := Phone{}
	camera := Camera{}
	Working(&phone)
	Working(&camera)

	// 1. 接口可以指向实现了该接口的变量
	// 2. 一个变量必须实现了接口的所有方法才能算是实现了接口
	// 3. 只要是自定义的数据类型，都可以实现接口
	// 4. 一个自定义数据类型可以实现多个接口
	var usb Usb = &Phone{}
	// 5. 多态通过接口实现，通过统一的接口来调用不同的实现，此时接口变量就呈现出多种形态
	usb.Start()

	// 类型断言的案例 多态数组
	var usbArr = [3]Usb{&Phone{}, &Phone{}, &Camera{}}
	// 遍历多态数组
	for _, usb := range usbArr {
		Working(usb)
	}
}

type AInterface interface {
	Test01()
}

type BInterface interface {
	Test02()
}

type CInterface interface {
	AInterface
	BInterface
	Test03()
}

type DStruct struct{}

func (d *DStruct) Test01() {
	fmt.Println("Test01 ...")
}

func (d *DStruct) Test02() {
	fmt.Println("Test02 ...")
}

func (d *DStruct) Test03() {
	fmt.Println("Test03 ...")
}

// 嵌套接口
func TestInterfaceNested(t *testing.T) {
	// 1. 若要实现一个继承了其他接口的接口，则必须要实现接口以及接口内嵌的接口的所有方法
	// 2. 若interface类型的变量没有指向任何实现了它的变量，那么它是引用类型，空接口
	// 3， 空接口interface{}没有任何方法，所有类型都实现了空接口，也就是说可以将任意类型的变量赋值给空接口类型的变量
	var c CInterface
	fmt.Printf("指向=%v, 类型=%T, 地址=%p\n", c, c, &c)
	c = &DStruct{}
	c.Test01()
	fmt.Printf("数据=%v, 类型=%T, 地址=%p\n", c, c, &c)
}

// 原结构体
type Hero struct {
	Name string
	Age  int
}

// 自定义类型，结构体切片
type HeroSlice []Hero

// 实现Len()方法，决定怎样获取长度
func (hs *HeroSlice) Len() int {
	return len(*hs)
}

// 实现Less()方法，决定使用什么标准进行排序
func (hs *HeroSlice) Less(i, j int) bool {
	// 若hs[i].Age > hs[j].Age表达式为true，则将hs[i]排到前面，否则排到后面(依次类推达成降序排列)
	return (*hs)[i].Age > (*hs)[j].Age
}

// 实现Swap()方法，决定怎样进行值的交换
func (hs *HeroSlice) Swap(i, j int) {
	(*hs)[i], (*hs)[j] = (*hs)[j], (*hs)[i]
}

// 自定义类型实现sort接口
func TestInterfaceSort(t *testing.T) {
	// 初始化随机数据
	var heroes HeroSlice
	for i := 0; i < 10; i++ {
		hero := Hero{
			Name: fmt.Sprintf("Hero-%d", i),
			Age:  rand.Intn(70),
		}
		heroes = append(heroes, hero)
	}

	// 随机数据
	for _, v := range heroes {
		fmt.Printf("%v\t", v)
	}
	fmt.Println()

	sort.Sort(&heroes)

	// 倒序排序
	for _, v := range heroes {
		fmt.Printf("%v\t", v)
	}
	fmt.Println()
}

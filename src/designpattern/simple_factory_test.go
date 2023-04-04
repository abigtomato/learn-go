package designpattern

import (
	"fmt"
	"testing"
)

// 水果的类型枚举
const (
	APPLE  = "apple"
	BANANA = "banana"
	PEAR   = "pear"
)

// 水果的抽象
type Fruit interface {
	Show()
}

// 水果的实现——苹果
type Apple struct{}

func (a *Apple) Show() {
	fmt.Println(APPLE)
}

// 水果的实现——香蕉
type Banana struct{}

func (b *Banana) Show() {
	fmt.Println(BANANA)
}

// 水果的实现——梨子
type Pear struct{}

func (p *Pear) Show() {
	fmt.Println(PEAR)
}

// 水果工厂
type FruitFactory struct{}

// 统一生产水果
func (f *FruitFactory) CreateFactory(name string) Fruit {
	switch name {
	case APPLE:
		return &Apple{}
	case BANANA:
		return &Banana{}
	case PEAR:
		return &Pear{}
	default:
		return nil
	}
}

// 简单工厂模式
func TestSimpleFactory(t *testing.T) {
	fruitFactory := &FruitFactory{}
	fruitFactory.CreateFactory(APPLE).Show()
	fruitFactory.CreateFactory(BANANA).Show()
	fruitFactory.CreateFactory(PEAR).Show()
}

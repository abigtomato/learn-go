package designpattern

import (
	"fmt"
	"testing"
)

// 抽象产品-苹果
type AbstractApple interface {
	ShowApple()
}

// 抽象产品-香蕉
type AbstractBanana interface {
	ShowBanana()
}

// 抽象产品-梨子
type AbstractPear interface {
	ShowPear()
}

// 抽象工厂
type AbstractFactory interface {
	CreateApple() AbstractApple
	CreateBanana() AbstractBanana
	CreatePear() AbstractPear
}

// 中国产品族：同一地区，不同功能的产品
type ChinaApple struct{}

func (ca *ChinaApple) ShowApple() {
	fmt.Println("ChinaApple")
}

// 中国产品族
type ChinaBanana struct{}

func (cb *ChinaBanana) ShowBanana() {
	fmt.Println("ChinaBanana")
}

// 中国产品族
type ChinaPear struct{}

func (cp *ChinaPear) ShowPear() {
	fmt.Println("ChinaPear")
}

// 中国产品族工厂：生成一组功能不同，地域相同的产品
type ChinaFactory struct{}

func (cf *ChinaFactory) CreateApple() AbstractApple {
	return new(ChinaApple)
}

func (cf *ChinaFactory) CreateBanana() AbstractBanana {
	return new(ChinaBanana)
}

func (cf *ChinaFactory) CreatePear() AbstractPear {
	return new(ChinaPear)
}

func Show(factory AbstractFactory) {
	factory.CreateApple().ShowApple()
	factory.CreateBanana().ShowBanana()
	factory.CreatePear().ShowPear()
}

// 抽象工厂模式
func TestAbstractFactory(t *testing.T) {
	Show(&ChinaFactory{})
}

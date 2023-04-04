package designpattern

import "testing"

// 水果工厂的抽象
type Factory interface {
	// 创建水果
	CreateFruit() Fruit
}

// 实现苹果工厂
type AppleFactory struct{}

// 苹果工厂只创建苹果
func (af *AppleFactory) CreateFruit() Fruit {
	return new(Apple)
}

// 实现香蕉工厂
type BananaFactory struct{}

// 香蕉工厂只创建香蕉
func (bf *BananaFactory) CreateFruit() Fruit {
	return new(Banana)
}

// 实现梨子工厂
type PearFactory struct{}

// 梨子工厂只创建梨子
func (pf *PearFactory) CreateFruit() Fruit {
	return new(Pear)
}

// 封装工厂的架构层
func FruitShow(factory Factory) {
	factory.CreateFruit().Show()
}

// 工厂方法模式
func TestFactoryMethod(t *testing.T) {
	FruitShow(new(AppleFactory))
	FruitShow(new(BananaFactory))
	FruitShow(new(PearFactory))
}

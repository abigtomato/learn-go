package designpattern

import (
	"fmt"
	"testing"
)

// 抽象策略-营销策略
type SellStrategy interface {
	GetPrice(price float64) float64
}

// 具体的策略A 打八折
type StrategyA struct{}

func (sa *StrategyA) GetPrice(price float64) float64 {
	return price * 0.8
}

// 具体的策略B 满200返现100
type StrategyB struct{}

func (sb *StrategyB) GetPrice(price float64) float64 {
	if price >= 200 {
		price -= 100
	}
	return price
}

// 环境类-商品，即策略的使用方
type SellGoods struct {
	Price    float64
	Strategy SellStrategy
}

func (sg *SellGoods) SetStrategy(s SellStrategy) {
	sg.Strategy = s
}

func (sg *SellGoods) SellPrice() float64 {
	fmt.Println("原价值", sg.Price, ".")
	return sg.Strategy.GetPrice(sg.Price)
}

// 策略模式
func TestStrategy(t *testing.T) {
	goods := SellGoods{Price: 200.0}

	goods.SetStrategy(new(StrategyA))
	fmt.Println(goods.SellPrice())

	goods.SetStrategy(new(StrategyB))
	fmt.Println(goods.SellPrice())
}

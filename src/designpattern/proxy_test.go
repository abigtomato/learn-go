package designpattern

import (
	"fmt"
	"testing"
)

// 商品
type Goods struct {
	Kind string
	Fact bool
}

// 抽象主题-购物
type Shopping interface {
	// 购买行为
	Buy(goods *Goods)
}

// 真实主题-在韩国购物
type KoreaShopping struct{}

func (ks *KoreaShopping) Buy(goods *Goods) {
	fmt.Println("KoreaShopping: " + goods.Kind)
}

// 真实主题-在美国购物
type AmericanShopping struct{}

func (as *AmericanShopping) Buy(goods *Goods) {
	fmt.Println("AmericanShopping: " + goods.Kind)
}

// 真实主题-在非洲购物
type AfricaShopping struct{}

func (as *AfricaShopping) Buy(goods *Goods) {
	fmt.Println("AfricaShopping: " + goods.Kind)
}

// 代理主题-海外代购
type OverseasProxy struct {
	// 用组合模式代理购物行为
	shopping Shopping
}

func (op *OverseasProxy) Buy(goods *Goods) {
	// 在购买之前植入了增强的行为
	if op.distinguish(goods) == true {
		op.shopping.Buy(goods)
		op.check(goods)
	}
}

// 代理增强行为-真假辨别
func (op *OverseasProxy) distinguish(goods *Goods) bool {
	fmt.Println("对[", goods.Kind, "]进行了辨别真伪.")
	if goods.Fact == false {
		fmt.Println("发现假货", goods.Kind, ", 不应该购买。")
	}
	return goods.Fact
}

// 代理增强行为-海关检测
func (op *OverseasProxy) check(goods *Goods) {
	fmt.Println("对[", goods.Kind, "] 进行了海关检查， 成功的带回祖国")
}

func NewShoppingProxy(shopping Shopping) Shopping {
	return &OverseasProxy{shopping}
}

// 代理模式
func TestProxy(t *testing.T) {
	NewShoppingProxy(new(AmericanShopping)).Buy(&Goods{
		Kind: "CET4",
		Fact: true,
	})
}

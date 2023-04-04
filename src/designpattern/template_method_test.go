package designpattern

import (
	"fmt"
	"testing"
)

// 饮料制作流程
type Beverage interface {
	BoilWater()          // 煮开水
	Brew()               // 冲泡
	PourInCup()          // 倒入杯中
	AddThings()          // 添加酌料
	WantAddThings() bool // 是否加入酌料Hook
}

// 饮料制作流程模版
type BeverageTemplate struct {
	b Beverage
}

// 固定的模板方法，封装一套固定的饮料制作流程
func (t *BeverageTemplate) MakeBeverage() {
	if t == nil {
		return
	}

	t.b.BoilWater()
	t.b.Brew()
	t.b.PourInCup()

	if t.b.WantAddThings() == true {
		t.b.AddThings()
	}
}

// 冲泡咖啡的流程
type MakeCaffe struct {
	// 组合冲泡饮料的流程模版
	t BeverageTemplate
}

func NewMakeCaffe() *MakeCaffe {
	makeCaffe := new(MakeCaffe)
	makeCaffe.t.b = makeCaffe
	return makeCaffe
}

func (mc MakeCaffe) BoilWater() {
	fmt.Println("将水煮到100摄氏度")
}

func (mc MakeCaffe) Brew() {
	fmt.Println("用水冲咖啡豆")
}

func (mc MakeCaffe) PourInCup() {
	fmt.Println("将充好的咖啡倒入陶瓷杯中")
}

func (mc MakeCaffe) AddThings() {
	fmt.Println("添加牛奶和糖")
}

func (mc MakeCaffe) WantAddThings() bool {
	return true
}

// 制作茶
type MakeTea struct {
	t BeverageTemplate
}

func NewMakeTea() *MakeTea {
	makeTea := new(MakeTea)
	makeTea.t.b = makeTea
	return makeTea
}

func (mt *MakeTea) BoilWater() {
	fmt.Println("将水煮到80摄氏度")
}

func (mt *MakeTea) Brew() {
	fmt.Println("用水冲茶叶")
}

func (mt *MakeTea) PourInCup() {
	fmt.Println("将充好的茶倒入茶壶中")
}

func (mt *MakeTea) AddThings() {
	fmt.Println("添加柠檬")
}

func (mt *MakeTea) WantAddThings() bool {
	return true
}

func TestTemplateMethod(t *testing.T) {
	NewMakeCaffe().t.MakeBeverage()
	NewMakeTea().t.MakeBeverage()
}

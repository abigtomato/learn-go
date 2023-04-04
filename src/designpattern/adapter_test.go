package designpattern

import (
	"fmt"
	"testing"
)

// V5-适配的目标
type V5 interface {
	Use5V()
}

// 业务类-计算机
type Computer struct {
	// 计算机只能使用V5充电
	v V5
}

func (c *Computer) Charge() {
	fmt.Println("Computer进行充电")
	c.v.Use5V()
}

func NewComputer(v V5) *Computer {
	return &Computer{v}
}

// V220-被适配的角色
type V220 struct{}

func (v *V220) Use220V() {
	fmt.Println("使用220V的电压")
}

// 适配器（用v5去适配v220）
type Adapter struct {
	// 组合被适配者
	v220 *V220
}

// 实现适配目标
func (a *Adapter) Use5V() {
	fmt.Println("使用适配器进行充电")
	// 内部使用被适配角色提供的方法
	a.v220.Use220V()
}

func NewAdapter(v220 *V220) *Adapter {
	return &Adapter{v220}
}

// 适配器模式
func TestAdapter(t *testing.T) {
	NewComputer(NewAdapter(new(V220))).Charge()
}

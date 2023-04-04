package designpattern

import (
	"fmt"
	"testing"
)

// 抽象的构件-手机
type Phone interface {
	Display()
}

// 抽象装饰器
type Decorator struct {
	phone Phone
}

func (d *Decorator) Display() {}

// 具体的构件-华为手机
type Huawei struct{}

// 构件具备的行为-显示
func (h *Huawei) Display() {
	fmt.Println("Huawei Display")
}

// 具体的构件-小米手机
type Xiaomi struct{}

func (x *Xiaomi) Display() {
	fmt.Println("Xiaomi Display")
}

// 具体的装饰器-手机膜
type MoDecorator struct {
	Decorator
}

func (md *MoDecorator) Display() {
	md.phone.Display()
	fmt.Println("MoDecorator Display")
}

func NewMoDecorator(phone Phone) Phone {
	return &MoDecorator{Decorator{phone}}
}

// 具体的装饰器-手机壳
type KeDecorator struct {
	Decorator
}

func (kd *KeDecorator) Display() {
	kd.phone.Display()
	fmt.Println("KeDecorator Display")
}

func NewKeDecorator(phone Phone) Phone {
	return &KeDecorator{Decorator{phone}}
}

// 装饰器模式
func TestDecorator(t *testing.T) {
	phone := new(Xiaomi)
	phone.Display()
	NewMoDecorator(phone).Display()
	NewKeDecorator(phone).Display()
}

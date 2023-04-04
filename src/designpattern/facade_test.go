package designpattern

import (
	"fmt"
	"testing"
)

// 子系统-电视
type TV struct{}

func (tv *TV) On() {
	fmt.Println("TV On")
}

func (tv *TV) Off() {
	fmt.Println("TV Off")
}

// 子系统-音响
type VoiceBox struct{}

func (vb *VoiceBox) On() {
	fmt.Println("VoiceBox On")
}

func (vb *VoiceBox) Off() {
	fmt.Println("VoiceBox Off")
}

// 子系统-灯光
type Light struct{}

func (light *Light) On() {
	fmt.Println("Light On")
}

func (light *Light) Off() {
	fmt.Println("Light Off")
}

// 子系统-游戏机
type Xbox struct{}

func (xbox *Xbox) On() {
	fmt.Println("Xbox On")
}

func (xbox *Xbox) Off() {
	fmt.Println("Xbox Off")
}

// 子系统-手机
type MicroPhone struct{}

func (mp *MicroPhone) On() {
	fmt.Println("MicroPhone On")
}

func (mp *MicroPhone) Off() {
	fmt.Println("MicroPhone Off")
}

type Projector struct{}

// 子系统-放映机
func (p *Projector) On() {
	fmt.Println("Projector On")
}

func (p *Projector) Off() {
	fmt.Println("Projector Off")
}

// 外观模式-提供家庭影院，组合各个模块
type HomePlayerFacade struct {
	tv    TV
	vb    VoiceBox
	light Light
	xbox  Xbox
	mp    MicroPhone
	pro   Projector
}

func (hp *HomePlayerFacade) DoKTV() {
	fmt.Println("家庭影院进入KTV模式")
	hp.tv.On()
	hp.pro.On()
	hp.mp.On()
	hp.light.Off()
	hp.vb.On()
}

func (hp *HomePlayerFacade) DoGame() {
	fmt.Println("家庭影院进入Game模式")
	hp.tv.On()
	hp.light.On()
	hp.xbox.On()
}

// 外观模式
func TestFacade(t *testing.T) {
	homePlayer := new(HomePlayerFacade)
	homePlayer.DoKTV()
	homePlayer.DoGame()
}

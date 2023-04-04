package designpattern

import (
	"fmt"
	"testing"
)

// 抽象的观察者
type Listener interface {
	// 观察某一类现象
	OnTeacherComing()
}

// 具体的观察者-学生（观察者，即被通知的目标）
type Student struct {
	BadThing string
}

func (s *Student) OnTeacherComing() {
	fmt.Println("张3 停止 ", s.BadThing)
}

func (s *Student) DoBadThing() {
	fmt.Println("张3 正在", s.BadThing)
}

// 抽象的通知者
type Notifier interface {
	AddListener(listener Listener)
	RemoveListener(listener Listener)
	Notify()
}

// 具体的通知者-班长（被观察者/目标/主题，即发送通知的一方）
type ClassMonitor struct {
	listenerList []Listener
}

func (cm *ClassMonitor) AddListener(listener Listener) {
	cm.listenerList = append(cm.listenerList, listener)
}

func (cm *ClassMonitor) RemoveListener(listener Listener) {
	for index, l := range cm.listenerList {
		if listener == l {
			cm.listenerList = append(cm.listenerList[:index], cm.listenerList[index+1:]...)
			break
		}
	}
}

func (cm *ClassMonitor) Notify() {
	for _, listener := range cm.listenerList {
		listener.OnTeacherComing()
	}
}

// 观察者模式
func TestObserver(t *testing.T) {
	zhang3 := &Student{BadThing: "吹牛逼"}
	zhao4 := &Student{BadThing: "玩游戏"}

	zhang3.DoBadThing()
	zhao4.DoBadThing()

	// 观察者
	classMonitor := new(ClassMonitor)
	classMonitor.AddListener(zhang3)
	classMonitor.AddListener(zhao4)

	// 通知被观察的目标
	classMonitor.Notify()
}

package designpattern

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// 懒汉式单例
type lazySingleton struct{}

// 全局唯一的实例，需要第一次使用时才会初始化
var lazyInstance *lazySingleton

func newLazyInstance() *lazySingleton {
	if lazyInstance != nil {
		return lazyInstance
	} else {
		lazyInstance = new(lazySingleton)
		return lazyInstance
	}
}

func (ls *lazySingleton) SomeThing() {
	fmt.Println("lazySingleton")
}

// 饿汉式单例
type hungrySingleton struct{}

// 包变量直接初始化
var hungryInstance = new(hungrySingleton)

func newHungryInstance() *hungrySingleton {
	return hungryInstance
}

func (hs *hungrySingleton) SomeThing() {
	fmt.Println("hungrySingleton")
}

// 线程安全的懒汉式单例
type syncLazySingleton struct{}

var syncLazyInstance *syncLazySingleton

// 标记
var initialized uint32

// 互斥锁
var lock sync.Mutex

func newSyncLazyInstance() *syncLazySingleton {
	// 优先走标记
	if atomic.LoadUint32(&initialized) == 1 {
		return syncLazyInstance
	}

	// 没有标记再加锁使用
	lock.Lock()
	defer lock.Unlock()

	if initialized == 0 {
		syncLazyInstance = new(syncLazySingleton)
		// 初始化完成后添加标记
		atomic.StoreUint32(&initialized, 1)
	}

	return syncLazyInstance
}

var once sync.Once

func newSyncLazyInstanceWithOnce() *syncLazySingleton {
	// 使用Go提供的once保证初始化操作只执行一次
	once.Do(func() {
		syncLazyInstance = new(syncLazySingleton)
	})
	return syncLazyInstance
}

func (hs *syncLazySingleton) SomeThing() {
	fmt.Println("syncLazyInstance")
}

// 单例模式
func TestSingleton(t *testing.T) {
	newLazyInstance().SomeThing()
	newHungryInstance().SomeThing()
	newSyncLazyInstance().SomeThing()
	newSyncLazyInstanceWithOnce().SomeThing()
}

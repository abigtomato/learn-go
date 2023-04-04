package goroutine

import (
	"fmt"
	"testing"
	"time"
)

// 管道基础语法
func TestChannelBasic(t *testing.T) {
	type Cat struct {
		Name string
		Age  int
	}

	// 创建一个 chan interface{} 空接口类型的管道（底层为队列结构）
	var iChan chan any

	// 无缓冲阻塞读写:
	// 1. 管道是引用类型，必须先通过make分配内存才可以使用
	// 2. 如果不定义缓冲区，则入队一个数据，就需要出队一个数据，否则阻塞读写
	// 有缓冲写满缓冲区后阻塞读写:
	// 1. 缓冲区长度4，定义了缓冲区才可以只入队不出队
	// 2. 等channel缓冲的数据满了，才会出现阻塞等待
	iChan = make(chan any, 4)
	fmt.Printf("指向=%v, 类型=%T, 地址=%p\n", iChan, iChan, &iChan)

	// iChan <- 向管道内部写入数据（入队操作）
	iChan <- make(map[string]int)
	iChan <- make([]float64, 10)
	iChan <- fmt.Sprintf("channel->%v", iChan)
	iChan <- &Cat{Name: "lily", Age: 3}
	fmt.Printf("获取缓冲区未读取的数据个数: %v, 获取缓冲区容量: %v\n", len(iChan), cap(iChan))

	// 使用内置 close() 函数关闭管道，使其只能读取数据无法写入数据
	close(iChan)

	// <-iChan 从管道取出数据（出队操作），channel 为空时再取数据则会报错
	mapper := <-iChan
	slice := <-iChan
	str := <-iChan
	cat := (<-iChan).(*Cat) // 通过类型断言转换为原始类型（因为存入channel时是以空接口类型存入的）
	fmt.Printf("map=%v, slice=%v, string=%v\n", mapper, slice, str)
	fmt.Printf("cat=%v, *cat.Name=%v\n", cat, (*cat).Name)

	// 管道的遍历测试
	intChan := make(chan int, 100)
	for i := 0; i < 100; i++ {
		intChan <- i * 2
	}
	// 遍历 channel 之前要保证关闭，使其不能再写入
	close(intChan)
	// 可以使用 range 取出 channel 中的数据，阻塞取出所有数据，否则 for 不会退出
	for v := range intChan {
		fmt.Printf("value: %v\n", v)
	}
}

// 只读只写管道
func TestWrChannel(t *testing.T) {
	// 数据和退出管道
	intChan := make(chan int, 10)
	exitChan := make(chan bool, 2)

	// chan<- 指定管道的状态为只写，适用于发送数据之类的场景
	go func(intChan chan<- int, exitChan chan bool) {
		defer close(intChan)

		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			intChan <- i
			fmt.Printf("发送数据: %v\n", i)
		}

		exitChan <- true
	}(intChan, exitChan)

	// <-chan 指定管道的状态为只读，适用于接收数据之类的场景
	go func(intChan <-chan int, exitChan chan bool) {
		for {
			if val, ok := <-intChan; !ok {
				break
			} else {
				fmt.Printf("接收数据: %v\n", val)
			}
		}

		exitChan <- true
	}(intChan, exitChan)

	// 主线程循环判断exitChan管道的结束标记，若存在2个结束标记则代表接收和发送协程都结束任务，之后主线程结束
	var total = 0
	for {
		if _, ok := <-exitChan; ok {
			total++
		}

		if total == 2 {
			break
		}
	}
}

// 生产者消费者模型:
// 1. 生产者: 发送数据端
// 2. 消费者: 接收数据端
// 3. 缓冲区:
// 3.1 解耦（降低生产者和消费者间的耦合度）
// 3.2 并发（生产者消费者数量不对等时，保持正常通信）
// 3.3 缓冲（生产者消费者处理速度不一致时，暂存数据）
func TestChannelProducerConsumer(t *testing.T) {
	// 管道做为模型中的缓冲区，无缓冲管道为模型提供同步通信，有缓冲管道则提供异步通信
	intChan := make(chan int, 50)
	exitChan := make(chan bool)

	// 生产者
	go func(intChan chan<- int) {
		// 生产完毕后关闭channel不再写入
		defer close(intChan)

		// intChan管道提供数据的生产消费
		for i := 0; i < 50; i++ {
			intChan <- i
			fmt.Printf("数据写入 -> %v\n", i)
			time.Sleep(time.Second)
		}
	}(intChan)

	// 消费者
	go func(intChan <-chan int, exitChan chan<- bool) {
		// 关闭标记管道让主go程判断是否消费结束
		defer close(exitChan)

		for {
			if v, ok := <-intChan; !ok {
				// 无缓冲管道关闭后再次读取会读出0
				// 有缓冲管道关闭后再次读取会先读出缓冲区的数据，读完后会读出0
				fmt.Printf("关闭后再次读取: %v\n", <-intChan)
				break
			} else {
				fmt.Printf("数据读取 -> %v\n", v)
			}
		}

		// exitChan管道提供消费结束的标识
		exitChan <- true
	}(intChan, exitChan)

	// 通过exitChan管道的标识判断主线程该何时结束
	for {
		// 若标记管道的写端关闭，再次读取的ok值为false
		if _, ok := <-exitChan; !ok {
			break
		}
	}
}

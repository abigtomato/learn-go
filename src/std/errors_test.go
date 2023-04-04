package std

import (
	"fmt"
	"testing"
	"time"
)

// 自定义错误的结构体
type CustomError struct {
	When time.Time // 什么时候发生错误
	What string    // 什么错误
}

// 结构体实现Error接口
func (e CustomError) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}

func oops() error {
	return CustomError{
		time.Date(1989, 3, 15, 22, 30, 0, 0, time.UTC),
		"the file system has gone away",
	}
}

func TestErrors(t *testing.T) {
	err := oops()
	if err != nil {
		fmt.Println(err)
	}
}

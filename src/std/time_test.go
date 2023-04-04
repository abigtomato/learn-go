package std

import (
	"fmt"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	// 基本使用
	now := time.Now()
	fmt.Println(now)
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)
}

func TestTimestamp(t *testing.T) {
	now := time.Now()
	// 时间戳是 1970年1月1日（08:00:00GTM）至今的总毫秒数，也称为 Unix 时间戳
	fmt.Println(now.Unix())
	fmt.Println(now.UnixNano())
}

func TestConv(t *testing.T) {
	// 将时间戳转为时间格式
	timestamp := time.Now().Unix()
	timeObj := time.Unix(timestamp, 0)
	fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n",
		timeObj.Year(), timeObj.Month(), timeObj.Day(),
		timeObj.Hour(), timeObj.Minute(), timeObj.Second())
}

func TestOps(t *testing.T) {
	// 时间操作
	now := time.Now()
	target := now.Add(time.Hour)
	// 时间增加
	fmt.Println(now.AddDate(1, 1, 1))
	// 时间间隔
	fmt.Println(target.Sub(now))
	// 相等判断
	fmt.Println(target.Equal(now))
	// 是否在之前
	fmt.Println(target.Before(now))
	// 是否在之后
	fmt.Println(target.After(now))
}

func TestFormat(t *testing.T) {
	now := time.Now()
	// Go 中格式化时间模板不是常见的 Y-m-d H:M:S，而是 go 诞生的时间2006年1月2号15点04分
	fmt.Println(now.Format("2006/01/02 15:04"))
	fmt.Println(now.Format("2006-01-02 15:04:05"))
	fmt.Println(now.Format("15:04 2006/01/02"))
}

func TestLocation(t *testing.T) {
	now := time.Now()
	// 加载时区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	// 按照指定时区指定格式解析字符串时间
	target, _ := time.ParseInLocation("2006/01/02 15:04:05", "2022/12/21 10:15:20", loc)
	fmt.Println(target, target.Sub(now))
}

func TestTick(t *testing.T) {
	// 使用 time.Tick(时间间隔) 来设置定时器，定时器本质是一个通道 channel
	ticker := time.Tick(time.Second)
	for i := range ticker {
		fmt.Println(i)
	}
}

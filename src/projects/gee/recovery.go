package gee

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 堆栈信息追踪
func trace(message string) string {
	var pcs [32]uintptr
	// Callers 用来返回调用栈的程序计数器
	// 第 0 个 Caller 是 Callers 本身，第 1 个是上一层 trace，第 2 个是再上一层的 defer func
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		// 获取对应的函数
		fn := runtime.FuncForPC(pc)
		// 获取到调用该函数的文件名和行号
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

// Recovery 异常捕获中间件
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, errors.New(message))
			}
		}()

		c.Next()
	}
}

package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]any

// Context 一次http请求的上下文
type Context struct {
	Writer     http.ResponseWriter // 响应抽象
	Req        *http.Request       // 请求抽象
	Path       string              // 路径
	Method     string              // 方法
	Params     map[string]string   // 参数
	StatusCode int                 // 状态码
	handlers   []HandlerFunc       // 处理器/中间件
	index      int                 // 当前执行到的中间件索引
	engine     *Engine             // 所属引擎
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	for s := len(c.handlers); c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...any) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	if _, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...))); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) JSON(code int, obj any) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	if err := json.NewEncoder(c.Writer).Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Fail(code int, err error) {
	http.Error(c.Writer, err.Error(), code)
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	if _, err := c.Writer.Write(data); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) HTML(code int, name string, data any) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(http.StatusInternalServerError, err)
	}
}

package gee

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/geektutu")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "geektutu" {
		t.Fatal("name should be equal to 'geektutu'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])
}

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := New()

	r.Use(Logger())
	r.Use(Recovery())
	r.GET("/panic", func(c *Context) {
		names := []string{"abigtomato"}
		c.String(http.StatusOK, names[100])
	})

	r = Default()

	//r.SetFuncMap(template.FuncMap{
	//	"FormatAsData": FormatAsDate,
	//})
	//r.LoadHTMLGlob("templates/*")
	//r.Static("/assets", "./static")

	//stu1 := &student{Name: "Geektutu", Age: 20}
	//stu2 := &student{Name: "Jack", Age: 22}
	//r.GET("/", func(c *gee.Context) {
	//	c.HTML(http.StatusOK, "css.tmpl", nil)
	//})
	//r.GET("/students", func(c *gee.Context) {
	//	c.HTML(http.StatusOK, "arr.tmpl", gee.H{
	//		"title":  "gee",
	//		"stuArr": [2]*student{stu1, stu2},
	//	})
	//})
	//r.GET("/date", func(c *gee.Context) {
	//	c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
	//		"title": "gee",
	//		"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
	//	})
	//})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *Context) {
			c.JSON(http.StatusOK, H{
				"name": "Hello Gee",
			})
		})

		v1.GET("/hello", func(c *Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.GET("/login", func(c *Context) {
			c.JSON(http.StatusOK, H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	_ = r.Run(":9999")
}

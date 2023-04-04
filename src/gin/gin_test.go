package gin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/errgroup"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"testing"
	"time"
)

// REST API
func TestREST(t *testing.T) {
	// Default 使用 Logger 和 Recovery 中间件
	r := gin.Default()
	// 不使用默认的中间件
	// r := gin.New()

	r.GET("/student", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "REST GET",
		})
	})

	r.POST("/create_student", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "REST POST",
		})
	})

	r.PUT("/update_student", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "REST PUT",
		})
	})

	r.DELETE("/delete_student", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "REST DELETE",
		})
	})

	// 自定义 HTTP 配置
	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()

	// 直接开启服务
	// _ = r.Run("localhost:8080")
}

// HTML渲染
func TestHTML(t *testing.T) {
	r := gin.Default()
	r.LoadHTMLFiles("./templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"name": "admin",
			"pwd":  "123456",
		})
	})

	_ = r.Run(":8080")
}

// 参数解析
func TestParams(t *testing.T) {
	r := gin.Default()
	r.LoadHTMLFiles("./templates/*")

	// GET /query_params?username=admin&password=123
	r.GET("/query_params", func(context *gin.Context) {
		username := context.Query("username")
		password := context.DefaultQuery("password", "123")
		context.JSON(http.StatusOK, gin.H{
			"name": username,
			"pwd":  password,
		})
	})

	// POST /form_params
	// Content-Type: application/x-www-form-urlencoded
	// username=admin&password=123
	r.POST("/form_params", func(context *gin.Context) {
		username := context.PostForm("username")
		password := context.DefaultPostForm("password", "123")
		context.JSON(http.StatusOK, gin.H{
			"name": username,
			"pwd":  password,
		})
	})

	type Person struct {
		Username string `uri:"id" binding:"required,uuid"`
		Password string `uri:"name" binding:"required"`
	}

	// uri绑定
	// GET /path_params/admin/123
	r.GET("/path_params/:username/:password", func(context *gin.Context) {
		//username := context.PostForm("username")
		//password := context.PostForm("password")
		var person Person
		if err := context.ShouldBindUri(&person); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		}
		context.JSON(http.StatusOK, gin.H{
			"name": person.Username,
			"pwd":  person.Password,
		})
	})

	type LoginForm struct {
		// 如果一个字段的 tag 加上了 binding:"required"，但绑定时是空值, Gin 会报错
		User     string `form:"user" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	// 模型绑定
	// $ curl -v --form user=user --form password=password http://localhost:8080/login
	r.POST("/login", func(c *gin.Context) {
		var form LoginForm
		// 使用 ShouldBind 绑定 multipart form
		// ShouldBind 会根据 Content-Type Header 推断使用哪个绑定器
		// 如果是 `GET` 请求，只使用 `Form` 绑定引擎（`query`）。
		// 如果是 `POST` 请求，首先检查 `content-type` 是否为 `JSON` 或 `XML`，然后再使用 `Form`（`form-data`）。
		if c.ShouldBind(&form) == nil {
			if form.User == "user" && form.Password == "password" {
				c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		}
	})

	r.Any("/testing", func(c *gin.Context) {
		var form LoginForm
		// 只绑定 url 查询参数而忽略 post 数据
		if c.ShouldBindQuery(&form) == nil {
			log.Println("====== Only Bind By Query String ======")
			log.Println(form.User)
			log.Println(form.Password)
		}
	})

	type FormA struct {
		Foo string `json:"foo" xml:"foo" binding:"required"`
	}

	type FormB struct {
		Bar string `json:"bar" xml:"bar" binding:"required"`
	}

	// 将 request body 绑定到不同的结构体中
	r.Any("/testing2", func(c *gin.Context) {
		var objA FormA
		var objB FormB
		// 读取 c.Request.Body 并将结果存入上下文。
		if errA := c.ShouldBindBodyWith(&objA, binding.JSON); errA == nil {
			c.String(http.StatusOK, `the body should be formA`)
		} else if errB := c.ShouldBindBodyWith(&objB, binding.JSON); errB == nil {
			// 这时, 复用存储在上下文中的 body
			c.String(http.StatusOK, `the body should be formB JSON`)
		} else if errB2 := c.ShouldBindBodyWith(&objB, binding.XML); errB2 == nil {
			// 可以接受其他格式
			c.String(http.StatusOK, `the body should be formB XML`)
		}
	})

	// 映射查询字符串或表单参数
	// POST /testing3?ids[a]=1234&ids[b]=hello
	// Content-Type: application/x-www-form-urlencoded
	// names[first]=thinkerou&names[second]=tianou
	r.POST("/testing3", func(c *gin.Context) {
		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")
		fmt.Printf("ids: %v; names: %v\n", ids, names)
	})

	_ = r.Run("localhost:8080")
}

// 参数检查
func TestParamCheck(t *testing.T) {
	// Booking 包含绑定和验证的数据
	type Booking struct {
		CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
		CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn,bookabledate" time_format:"2006-01-02"`
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义验证器
		_ = v.RegisterValidation("bookabledate", func(fl validator.FieldLevel) bool {
			if date, ok := fl.Field().Interface().(time.Time); ok {
				today := time.Now()
				if today.After(date) {
					return false
				}
			}
			return true
		})
	}

	route := gin.Default()

	route.GET("/bookable", func(c *gin.Context) {
		var b Booking
		if err := c.ShouldBindWith(&b, binding.Query); err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	_ = route.Run(":8080")
}

// 路由
func TestRoute(t *testing.T) {
	engine := gin.Default()

	// 路由组
	userGroup := engine.Group("/user")
	{
		userGroup.GET("/index", func(context *gin.Context) {})
		userGroup.POST("/login", func(context *gin.Context) {})
	}

	// 匹配任意请求
	engine.Any("/test", func(context *gin.Context) {})

	// 未匹配到路由的请求
	engine.NoRoute(func(context *gin.Context) {})

	// 此 handler 将匹配 /user/john 但不会匹配 /user/ 或者 /user
	engine.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// 此 handler 将匹配 /user/john/ 和 /user/john/send
	// 如果没有其他路由匹配 /user/john，它将重定向到 /user/john/
	engine.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	_ = engine.Run()
}

// 中间件、处理器、钩子函数
func TestHook(t *testing.T) {
	engine := gin.Default()

	// 自定义中间件
	engine.Use(func(c *gin.Context) {
		now := time.Now()

		c.Set("example", "123456")

		// 请求前

		c.Next()

		// 请求后

		latency := time.Since(now)
		log.Print(latency)

		// 获取发送的status
		status := c.Writer.Status()
		log.Println(status)
	})

	// LoggerWithFormatter 中间件会写入日志到 gin.DefaultWriter
	// 默认 gin.DefaultWriter = os.Stdout
	engine.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		// 返回自定义格式
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			params.ClientIP,
			params.TimeStamp.Format(time.RFC1123),
			params.Method,
			params.Path,
			params.Request.Proto,
			params.StatusCode,
			params.Latency,
			params.Request.UserAgent(),
			params.ErrorMessage,
		)
	}))

	// 从任何 panic 中恢复，并响应 500
	engine.Use(gin.Recovery())

	engine.GET("/", func(context *gin.Context) {
		// 中间件、处理器、钩子函数
		// 在中间件或 handler 中启动新的 Goroutine 时，不能使用原始的上下文，必须使用只读副本
		tmp := context.Copy()
		go func() {
			time.Sleep(5 * time.Second)
			log.Println("Done! in path " + tmp.Request.URL.Path)
		}()
	}, func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	})

	_ = engine.Run()
}

// 结构体绑定
func TestStructBind(t *testing.T) {
	type UserInfo struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}

	r := gin.Default()

	r.GET("/user", func(c *gin.Context) {
		var user UserInfo
		// 基于请求自动提取JSON，Form表单，Query等类型的值，并把值绑定到指定的结构体对象
		// 如果是GET请求，只使用Form绑定引擎（Query）
		// 如果是POST请求，首先检查content-type是否为JSON或XML，然后再使用Form（form-data）
		err := c.ShouldBind(&user)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		}
		fmt.Printf("%#v\v", user)
	})

	_ = r.Run()
}

// 文件上传
func TestUpload(t *testing.T) {
	r := gin.Default()

	// 处理multipart/form-data提交文件时默认的内存限制是32 MiB
	// 可以通过下面的方式修改
	// r.MaxMultipartMemory = 8 << 20  // 8 MiB
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("f1")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		dst := fmt.Sprintf("D:/tmp/%s", file.Filename)
		_ = c.SaveUploadedFile(file, dst)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	})

	// 多文件上传
	r.POST("/multi_upload", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["file"]
		for _, file := range files {
			dst := fmt.Sprintf("D:/tmp/%s", file.Filename)
			_ = c.SaveUploadedFile(file, dst)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%d files uploaded!", len(files)),
		})
	})

	_ = r.Run()
}

// 重定向
func TestRedirect(t *testing.T) {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		// http重定向
		c.Redirect(http.StatusPermanentRedirect, "https://www.baidu.com")
	})

	r.GET("/test1", func(c *gin.Context) {
		// 路由重定向
		c.Request.URL.Path = "/test2"
		r.HandleContext(c)
	})
	r.GET("/test2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Hello": "Golang",
		})
	})

	_ = r.Run()
}

// Cookie操作
func TestCookie(t *testing.T) {
	r := gin.Default()

	r.GET("/cookie", func(c *gin.Context) {
		// 获取cookie
		cookie, err := c.Cookie("gin_cookie")
		if err != nil {
			cookie = "NotStr"
			// 设置cookie
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}
		fmt.Printf("Cookie value: %s\n", cookie)
	})

	_ = r.Run()
}

// 日志操作
func TestLog(t *testing.T) {
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色
	gin.DisableConsoleColor()
	// 强制日志颜色化
	// gin.ForceConsoleColor()

	// 如果想要以指定的格式（例如 JSON，key values 或其他格式）记录信息，则可以使用 gin.DebugPrintRouteFunc 指定格式
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	// 将日志写入文件和控制台
	file, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	_ = r.Run()
}

// HTML绑定
func TestHTMLBind(t *testing.T) {
	type MyForm struct {
		Colors []string `form:"colors[]"`
	}

	r := gin.Default()

	r.LoadHTMLGlob("./templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "checkbox.html", nil)
	})
	r.POST("/", func(c *gin.Context) {
		var fakeForm MyForm
		// 和html对应的name绑定
		_ = c.Bind(&fakeForm)
		c.JSON(200, gin.H{
			"color": fakeForm.Colors,
		})
	})

	_ = r.Run()
}

// AsciiJson
func TestAsciiJSON(t *testing.T) {
	r := gin.Default()

	r.GET("/json", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "golang",
			"tag":  "<br>",
		}
		// 使用 AsciiJSON生成具有转义的非 ASCII字符的 ASCII-only JSON
		// 输出 : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.AsciiJSON(http.StatusOK, data)
	})

	_ = r.Run()
}

// 服务器推送
func TestPusher(t *testing.T) {
	var html = template.Must(template.New("https").Parse(`
	<html>
	<head>
	  <title>Https Test</title>
	  <script src="/assets/app.js"></script>
	</head>
	<body>
	  <h1 style="color:red;">Welcome, Ginner!</h1>
	</body>
	</html>
	`))

	r := gin.Default()

	r.Static("/assets", "./assets")
	r.SetHTMLTemplate(html)

	r.GET("/", func(c *gin.Context) {
		// HTTP2 server 推送
		if pusher := c.Writer.Pusher(); pusher != nil {
			if err := pusher.Push("/assets/app.js", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
		c.HTML(200, "https", gin.H{
			"status": "Susses",
		})
	})

	_ = r.Run(":8080", "./testdata/server.pem", "./testdata/server.key")
}

// JSONP
func TestJSONP(t *testing.T) {
	r := gin.Default()

	r.GET("/jsonp", func(c *gin.Context) {
		data := map[string]interface{}{
			"foo": "bar",
		}
		// 使用 JSONP向不同域的服务器请求数据。如果查询参数存在回调，则将回调添加到响应体中。
		// /JSONP?callback=x
		// 将输出：x({\"foo\":\"bar\"})
		c.JSONP(http.StatusOK, data)
	})

	_ = r.Run()
}

// PureJSON
func TestPureJSON(t *testing.T) {
	r := gin.Default()

	r.GET("/purseJson", func(c *gin.Context) {
		// JSON使用 unicode替换特殊 HTML字符，例如 < 变为 \ u003c。如果要按字面对这些字符进行编码，则可以使用 PureJSON
		c.PureJSON(200, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})

	_ = r.Run()
}

// SecureJSON
func TestSecureJSON(t *testing.T) {
	r := gin.Default()

	r.GET("/secureJson", func(c *gin.Context) {
		names := []string{"lena", "austin", "foo"}
		// 使用 SecureJSON防止 json劫持。如果给定的结构是数组值，则默认预置 "while(1)," 到响应体
		// 将输出：while(1);["lena","austin","foo"]
		c.SecureJSON(http.StatusOK, names)
	})

	_ = r.Run()
}

// JSON\XML\YAML\Protobuf渲染
func TestRendering(t *testing.T) {
	r := gin.Default()

	r.GET("/someJson", func(c *gin.Context) {
		// gin.H 是 map[string]interface{} 的一种快捷方式
		c.JSON(http.StatusOK, gin.H{
			"message": "hey",
			"status":  http.StatusOK,
		})
	})

	r.GET("/moreJson", func(c *gin.Context) {
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		// 将输出：{"user": "Lena", "Message": "hey", "Number": 123}
		c.JSON(http.StatusOK, msg)
	})

	r.GET("/someXml", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/someYaml", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/someProtoBuf", func(c *gin.Context) {
		reps := []int64{int64(1), int64(2)}
		label := "test"

		// protobuf 的具体定义写在 testdata/protoexample 文件中
		data := &protoexample.Test{
			Label: &label,
			Reps:  reps,
		}

		// 请注意，数据在响应中变为二进制数据
		// 将输出被 protoexample.Test protobuf 序列化了的数据
		c.ProtoBuf(http.StatusOK, data)
	})

	_ = r.Run()
}

// 优雅关闭
func TestShutdown(t *testing.T) {
	engine := gin.Default()
	engine.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Println("Server exiting")
}

// 从 reader 读取数据
func TestReader(t *testing.T) {
	router := gin.Default()

	router.GET("/someDataFromReader", func(c *gin.Context) {
		resp, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
		if err != nil || resp.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := resp.Body
		contentLength := resp.ContentLength
		contentType := resp.Header.Get("Content-Type")

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})

	_ = router.Run(":8080")
}

// 使用 BasicAuth 中间件
func TestBasicAuth(t *testing.T) {
	secrets := gin.H{
		"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
		"austin": gin.H{"email": "austin@example.com", "phone": "666"},
		"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
	}

	router := gin.Default()

	// 路由组使用 gin.BasicAuth() 中间件
	// gin.Accounts 是 map[string]string 的一种快捷方式
	authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	authorized.GET("/secrets", func(c *gin.Context) {
		// 获取用户，它是由 BasicAuth 中间件设置的
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	_ = router.Run(":8080")
}

// 全局中间件
func TestMiddleware(t *testing.T) {
	// 新建一个没有任何默认中间件的路由
	r := gin.New()

	// Logger 中间件将日志写入 gin.DefaultWriter，即使你将 GIN_MODE 设置为 release。
	r.Use(gin.Logger())

	// Recovery 中间件会 recover 任何 panic。如果有 panic 的话，会写入 500。
	r.Use(gin.Recovery())

	_ = r.Run(":8080")
}

// 运行多个服务
func TestManyService(t *testing.T) {
	var g errgroup.Group

	g.Go(func() error {
		server01 := &http.Server{
			Addr: ":8080",
			Handler: func() http.Handler {
				engine := gin.New()
				engine.Use(gin.Recovery())
				engine.GET("/", func(c *gin.Context) {
					c.JSON(
						http.StatusOK,
						gin.H{
							"code":  http.StatusOK,
							"error": "Welcome server 01",
						},
					)
				})
				return engine
			}(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		return server01.ListenAndServe()
	})

	g.Go(func() error {
		server02 := &http.Server{
			Addr: ":8081",
			Handler: func() http.Handler {
				engine := gin.New()
				engine.Use(gin.Recovery())
				engine.GET("/", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code":  http.StatusOK,
						"error": "Welcome server 02",
					})
				})
				return engine
			}(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		return server02.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

// 静态文件服务
func TestStaticServer(t *testing.T) {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.StaticFS("/more_static", http.Dir("my_file_system"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	_ = router.Run(":8080")
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

// 编写 Gin 的测试用例
func TestPingRoute(t *testing.T) {
	router := SetupRouter()

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "pong", recorder.Body.String())
}

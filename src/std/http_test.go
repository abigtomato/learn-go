package std

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	// GET请求
	r, err := http.Get("https://apis.juhe.cn/simpleWeather/query?key=087d7d10f700d20e27bb753cd806e40b&city=北京")
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
	b, _ := io.ReadAll(r.Body)
	fmt.Printf("b: %v\n", string(b))
}

func TestGetParam(t *testing.T) {
	params := url.Values{}
	Url, err := url.Parse("https://apis.juhe.cn/simpleWeather/query")
	if err != nil {
		return
	}
	params.Set("key", "087d7d10f700d20e27bb753cd806e40b")
	params.Set("city", "北京")
	// 如果参数中有中文参数，这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath)
	resp, err := http.Get(urlPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func TestGetJson(t *testing.T) {
	type result struct {
		Args    string            `json:"args"`
		Headers map[string]string `json:"headers"`
		Origin  string            `json:"origin"`
		Url     string            `json:"url"`
	}
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	var res result
	_ = json.Unmarshal(body, &res)
	fmt.Printf("%#v", res)
}

func TestGetHead(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://httpbin.org/get", nil)
	req.Header.Add("name", "zs")
	req.Header.Add("age", "80")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf(string(body))
}

func TestPost(t *testing.T) {
	path := "https://apis.juhe.cn/simpleWeather/query"
	urlValues := url.Values{}
	urlValues.Add("key", "087d7d10f700d20e27bb753cd806e40b")
	urlValues.Add("city", "北京")
	r, err := http.PostForm(path, urlValues)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
	b, _ := io.ReadAll(r.Body)
	fmt.Printf("b: %v\n", string(b))
}

func TestPostBody(t *testing.T) {
	urlValues := url.Values{
		"name": {"zs"},
		"age":  {"80"},
	}
	reqBody := urlValues.Encode()
	resp, _ := http.Post("https://httpbin.org/post",
		"text/html",
		strings.NewReader(reqBody))
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func TestPostJson(t *testing.T) {
	data := make(map[string]interface{})
	data["site"] = "www.duoke360.com"
	data["name"] = "多课网"
	bytesData, _ := json.Marshal(data)
	resp, _ := http.Post("https://httpbin.org/post",
		"application/json",
		bytes.NewReader(bytesData))
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func TestClient(t *testing.T) {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest(http.MethodGet,
		"https://apis.juhe.cn/simpleWeather/query?key=087d7d10f700d20e27bb753cd806e40b&city=北京",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("referer", "https://apis.juhe.cn/")
	res, err2 := client.Do(req)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	b, _ := io.ReadAll(res.Body)
	fmt.Printf("b: %v\n", string(b))
}

func TestServer(t *testing.T) {
	// 请求处理函数
	f := func(resp http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(resp, "hello world")
	}
	// 响应路径,注意前面要有斜杠 /
	http.HandleFunc("/hello", f)
	// 设置监听端口，并监听，注意前面要有冒号:
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatal(err)
	}
}

package testing

import (
	"io"
	"net/http/httptest"
	"testing"
)

func TestConn(t *testing.T) {
	// 使用 httptest 模拟请求对象(req)和响应对象(w)
	req := httptest.NewRequest("GET", "https://example.com/hello", nil)
	w := httptest.NewRecorder()
	HelloHandler(w, req)
	bytes, _ := io.ReadAll(w.Result().Body)

	if string(bytes) != "hello world" {
		t.Fatal("expected hello world, but got", string(bytes))
	}
}

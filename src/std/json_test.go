package std

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Person struct {
	Name  string
	Age   int
	Email string
}

func TestJson(t *testing.T) {
	p := Person{
		Name:  "tom",
		Age:   20,
		Email: "tom@gmail.com",
	}

	// 结构体转json
	b, _ := json.Marshal(p)
	fmt.Printf("string(b): %v\n", string(b))

	// json转结构体
	b1 := []byte(`{"Name":"tom","Age":20,"Email":"tom@gmail.com"}`)
	var p2 Person
	_ = json.Unmarshal(b1, &p2)
	fmt.Printf("p2: %v\n", p2)

	// 解析嵌套类型
	b2 := []byte(`{"Name":"tom","Age":20,"Email":"tom@gmail.com","Parents":["big tom","kite"]}`)
	var f map[string]interface{}
	_ = json.Unmarshal(b2, &f)
	fmt.Printf("f: %v\n", f)

	// 解析嵌套引用类型
	type Person2 struct {
		Name   string
		Age    int
		Email  string
		Parent []string
	}
	p3 := Person2{
		Name:   "tom",
		Age:    18,
		Email:  "tom@gmail.com",
		Parent: []string{"big tom", "big kite"},
	}
	b, _ = json.Marshal(p3)
	fmt.Printf("b: %v\n", b)
	fmt.Printf("string(b): %v\n", string(b))
}

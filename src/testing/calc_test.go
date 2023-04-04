package testing

import (
	"fmt"
	"os"
	"testing"
)

// 1. 测试用例名称一般命名为 Test 加上待测试的方法名
// 2. 测试用的参数有且只有一个，在这里是 t *testing.T
// 3. 基准测试benchmark的参数是 *testing.B，TestMain 的参数是 *testing.M 类型
// 4. go test -v，-v 参数会显示每个用例的测试结果，另外 -cover 参数可以查看覆盖率
// 5. 只运行其中的一个用例，例如 TestAdd，可以用 -run 参数指定，该参数支持通配符 *，和部分正则表达式，例如 ^、$
func TestAdd(t *testing.T) {
	if ans := Add(1, 2); ans != 3 {
		// t.Errorf 遇错不停
		t.Errorf("1 + 2 expected be 3, but %d got", ans)
	}
}

func TestMul(t *testing.T) {
	cases := []struct {
		Name           string
		A, B, Expected int
	}{
		{"", 2, 3, 6},
		{"", 2, -3, -6},
		{"", 2, 0, 0},
	}

	for _, c := range cases {
		// 子测试
		t.Run(c.Name, func(t *testing.T) {
			if ans := Mul(c.A, c.B); ans != c.Expected {
				// t.Fatal 遇错即停
				t.Fatalf("%d + %d expected %d, but %d got", c.A, c.B, c.Expected, ans)
			}
		})
	}
}

type calcCase struct {
	A, B, Expected int
}

// 对一些重复的逻辑，抽取出来作为公共的帮助函数helpers，可以增加测试代码的可读性和可维护性
func createMulTestCase(t *testing.T, c *calcCase) {
	// Go 语言在 1.9 版本中引入了 t.Helper()，用于标注该函数是帮助函数，报错时将输出帮助函数调用者的信息，而不是帮助函数的内部信息
	t.Helper()
	if ans := Mul(c.A, c.B); ans != c.Expected {
		t.Fatalf("%d * %d expected %d, but %d got", c.A, c.B, c.Expected, ans)
	}
}

func TestMul2(t *testing.T) {
	createMulTestCase(t, &calcCase{2, 3, 6})
	createMulTestCase(t, &calcCase{2, -3, -6})
	createMulTestCase(t, &calcCase{2, 0, 1})
}

func setup() {
	fmt.Println("Before all tests")
}

func teardown() {
	fmt.Println("After all tests")
}

// 如果测试文件中包含函数 TestMain，那么生成的测试将调用 TestMain(m)，而不是直接运行测试
func TestMain(m *testing.M) {
	// 额外的准备
	setup()
	// 调用 m.Run() 触发所有测试用例的执行
	code := m.Run()
	// 回收工作
	teardown()
	// 使用 os.Exit() 处理返回的状态码，如果不为0，说明有用例失败
	os.Exit(code)
}

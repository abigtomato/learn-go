package service

import (
	"go/ast"
	"log"
	"reflect"
	"sync/atomic"
)

// MethodType RPC服务方法类型抽象
// func (t *T) MethodName(argType T1, replyType *T2) error
// 第一个参数是入参，第二个参数是指针，即返回值
type MethodType struct {
	method    reflect.Method // 方法自身
	ArgType   reflect.Type   // 参数类型（第一个参数类型）
	ReplyType reflect.Type   // 返回值类型（第二个参数类型）
	numCalls  uint64         // 调用次数
}

func (m *MethodType) NumCalls() uint64 {
	return atomic.LoadUint64(&m.numCalls)
}

// NewArgv 创建入参对应类型的实例
func (m *MethodType) NewArgv() reflect.Value {
	var argv reflect.Value
	if m.ArgType.Kind() == reflect.Pointer {
		// 创建指针类型
		argv = reflect.New(m.ArgType.Elem())
	} else {
		// 创建值类型
		argv = reflect.New(m.ArgType).Elem()
	}
	return argv
}

// NewReplyv 创建返回值对应类型的实例（返回值必须是指针类型）
func (m *MethodType) NewReplyv() reflect.Value {
	replyv := reflect.New(m.ReplyType.Elem())
	switch m.ReplyType.Elem().Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}
	return replyv
}

// Service RPC服务的抽象
type Service struct {
	Name   string                 // 服务名 即映射的结构体的名称
	Typ    reflect.Type           // 服务类型 即映射的结构体的类型
	Rcvr   reflect.Value          // 结构体实例本身
	Method map[string]*MethodType // 存储映射的结构体的所有符合条件的方法
}

// NewService 映射结构体为服务
func NewService(rcvr any) *Service {
	s := new(Service)
	s.Rcvr = reflect.ValueOf(rcvr)
	s.Name = reflect.Indirect(s.Rcvr).Type().Name()
	s.Typ = reflect.TypeOf(rcvr)
	if !ast.IsExported(s.Name) {
		log.Fatalf("rpc server: %s is not a valid service name", s.Name)
	}
	s.registerMethods()
	return s
}

// 注册结构体方法为服务的方法
func (s *Service) registerMethods() {
	s.Method = make(map[string]*MethodType)
	for i := 0; i < s.Typ.NumMethod(); i++ {
		method := s.Typ.Method(i)
		mType := method.Type
		// 参数数量为3 第0个是自身，第1个是入参，第2个是指向返回值的指针
		// 返回值数量为1 类型为error
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		// 返回值类型必须是error
		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		// 参数类型必须是导出或内置类型
		argType, replyType := mType.In(1), mType.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}
		s.Method[method.Name] = &MethodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
		log.Printf("rpc server: register %s.%s\n", s.Name, method.Name)
	}
}

// 判断是否是导出或内置类型
func isExportedOrBuiltinType(t reflect.Type) bool {
	// 导出类型：类型名称以大写字母开头
	// 内置类型：包路径为空
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

// Call 方法调用
func (s *Service) Call(m *MethodType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)
	f := m.method.Func
	returnValues := f.Call([]reflect.Value{s.Rcvr, argv, replyv})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}

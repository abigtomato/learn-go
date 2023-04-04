package session

import (
	"Golearn/src/projects/geeorm/log"
	"reflect"
)

// 构造函数的扩展点
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

func (s *Session) CallMethod(method string, value any) {
	// 获取当前会话正在操作的数据表对象，并通过反射获取其指定方法
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		// 如果指定对象，则获取指定对象的方法
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	// 构造入参（钩子函数的入参类型都是*Session）
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		// 反射调用钩子函数
		if v := fm.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}

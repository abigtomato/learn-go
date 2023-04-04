package schema

import (
	"go/ast"
	"learn-go/src/projects/geeorm/dialect"
	"reflect"
)

// Field 表示数据库的字段列
type Field struct {
	Name string // 字段名
	Type string // 类型
	Tag  string // 约束条件
}

// Schema 表示数据库表结构
type Schema struct {
	Model      any               // 被映射的对象
	Name       string            // 表名
	Fields     []*Field          // 字段列表
	FieldNames []string          // 列名列表
	fieldMap   map[string]*Field // 记录结构体字段和数据表列的映射
}

// GetField 根据名称获取指定字段
func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

// RecordValues 根据数据库的列，从对象中获取对应值，按顺序平铺
func (s *Schema) RecordValues(dest any) []any {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []any
	for _, field := range s.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

// Parse 将go类型根据方言转换为数据库类型
func Parse(dest any, d dialect.Dialect) *Schema {
	// reflect.ValueOf反射获取dest的类型
	// reflect.Indirect获取指针指向的实例
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(), // 结构体名称作为表名
		fieldMap: make(map[string]*Field),
	}

	// NumField获取实例的字段个数
	for i := 0; i < modelType.NumField(); i++ {
		// 通过下标获取到具体字段
		sf := modelType.Field(i)
		if !sf.Anonymous && ast.IsExported(sf.Name) {
			field := &Field{
				Name: sf.Name,
				// 转换为数据库对应的字段
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(sf.Type))),
			}
			// 处理tag中标注的额外约束
			if v, ok := sf.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, sf.Name)
			schema.fieldMap[sf.Name] = field
		}
	}
	return schema
}

package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

// Dialect 数据库方言
type Dialect interface {
	DataTypeOf(typ reflect.Value) string            // 将go类型转换为特定数据库类型
	TableExistSQL(tableName string) (string, []any) // 返回某个表是否存在的 SQL 语句
}

// RegisterDialect 方言注册
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect 根据数据库获取方言
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}

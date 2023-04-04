package clause

import (
	"fmt"
	"strings"
)

// sql子句生成器
type generator func(values ...any) (string, []any)

// 子句类型和生成器映射
var generators map[Type]generator

// 初始化各个子句的生成器
func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

// INSERT INTO $tableName ($fields)
func _insert(values ...any) (string, []any) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []any{}
}

// VALUES ($v1), ($v2), ...
func _values(values ...any) (string, []any) {
	var bindStr string
	var sql strings.Builder
	var vars []any
	sql.WriteString("VALUES ")

	for i, value := range values {
		v := value.([]any)
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))

		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

// SELECT $fields FROM $tableName
func _select(values ...any) (string, []any) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []any{}
}

// LIMIT $num
func _limit(values ...any) (string, []any) {
	return "LIMIT ?", values
}

// WHERE $desc
func _where(values ...any) (string, []any) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

// ORDER BY $fields [DESC|ASC]
func _orderBy(values ...any) (string, []any) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []any{}
}

// UPDATE $tableName SET $fields
func _update(values ...any) (string, []any) {
	tableName := values[0]
	m := values[1].(map[string]any)
	var keys []string
	var vars []any
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

// DELETE FROM $tableName
func _delete(values ...any) (string, []any) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []any{}
}

// SELECT count(*) FROM $tableName
func _count(values ...any) (string, []any) {
	return _select(values[0], []string{"count(*)"})
}

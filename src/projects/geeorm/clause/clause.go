package clause

import "strings"

// Clause sql子句的封装
type Clause struct {
	sql     map[Type]string // 子句类型和sql语句的映射
	sqlVars map[Type][]any  // 子句类型和sql对应参数的映射
}

type Type int

// 子句类型
const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

// Set 根据子句类型生成对应的sql（调用生成器）
func (c *Clause) Set(name Type, vars ...any) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]any)
	}
	// 调用对应的子句生成器获取sql
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// Build 根据传入的子句类型，按顺序构造出最终的sql语句
func (c *Clause) Build(orders ...Type) (string, []any) {
	var sqlSlice []string
	var vars []any
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqlSlice = append(sqlSlice, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqlSlice, " "), vars
}

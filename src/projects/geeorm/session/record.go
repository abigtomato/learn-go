package session

import (
	"errors"
	"learn-go/src/projects/geeorm/clause"
	"reflect"
)

// Insert 将传入的对象字段平铺，构造成插入语句
func (s *Session) Insert(values ...any) (int64, error) {
	s.CallMethod(BeforeInsert, nil)

	recordValues := make([]any, 0)
	for _, value := range values {
		// 解析对象转换成表结构
		table := s.Model(value).RefTable()
		// 构造Insert子句
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		// 将传入的对象字段平铺
		recordValues = append(recordValues, table.RecordValues(value))
	}

	// 构造values子句
	s.clause.Set(clause.VALUES, recordValues...)
	// 生成sql
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)

	// 执行sql
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	s.CallMethod(AfterInsert, nil)

	// 返回受影响的行数
	return result.RowsAffected()
}

// Find 根据查询出的记录构造结果对象切片
func (s *Session) Find(values any) error {
	s.CallMethod(BeforeQuery, nil)

	// 指向结果切片的指针
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	// 获取切片的单个元素类型
	destType := destSlice.Type().Elem()
	// 通过元素类型构造出对应表结构
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	// 构造select子句
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)

	// 执行查询
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		// 反射创建结果对象实例
		dest := reflect.New(destType).Elem()
		var values []any
		// 将结果对象的字段平铺开
		for _, name := range table.FieldNames {
			// Addr是将地址存入values，这样在调用Scan时会将值真正赋值给结果对象而不是副本
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		// 将数据表查询的该行记录依次赋值给结果对象的每个字段
		if err := rows.Scan(values...); err != nil {
			return err
		}

		s.CallMethod(AfterQuery, dest.Addr().Interface())

		// 将赋值后的结果对象追加到结果切片中
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

// Update 接受2种入参，平铺开来的键值对和map类型的键值对
func (s *Session) Update(kv ...any) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)

	m, ok := kv[0].(map[string]any)
	if !ok {
		// 若是切片类型的平铺键值对，则转换为map类型
		m = make(map[string]any)
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}

	// 构造 UPDATE ... WHERE ... 子句并执行
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	s.CallMethod(AfterUpdate, nil)

	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)

	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	s.CallMethod(AfterDelete, nil)

	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...any) *Session {
	var vars []any
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

// First 返回一条记录
func (s *Session) First(value any) error {
	// 获取结果对象的指针
	dest := reflect.Indirect(reflect.ValueOf(value))
	// 反射创建结果对象的切片
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	// 构造limit 1的sql
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	// 返回第一条
	dest.Set(destSlice.Index(0))
	return nil
}

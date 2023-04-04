package geeorm

import (
	"Golearn/src/projects/geeorm/dialect"
	"Golearn/src/projects/geeorm/log"
	"Golearn/src/projects/geeorm/session"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

// Engine ORM引擎
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// NewSession 创建会话
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

type TxFunc func(*session.Session) (any, error)

// Transaction 事务的一键式使用接口
func (e *Engine) Transaction(f TxFunc) (result any, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}

	// 如果事务操作出现错误，立即回滚，否则提交
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}
	}()

	// 用户需要被事务包裹的操作通过回调函数封装
	return f(s)
}

// 用来计算前后两个字段切片的差集（新表 - 旧表 = 新增字段，旧表 - 新表 = 删除字段）
func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate 数据库迁移（结构体字段发送变化时变更数据表）
func (e *Engine) Migrate(value any) error {
	_, err := e.Transaction(func(s *session.Session) (result any, err error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}

		// 计算新旧结构的字段差异
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)

		// 处理新字段
		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}

		// 根据最新的字段创建临时表 删除旧表 重命名新表
		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}

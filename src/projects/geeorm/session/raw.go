package session

import (
	"database/sql"
	"learn-go/src/projects/geeorm/clause"
	"learn-go/src/projects/geeorm/dialect"
	"learn-go/src/projects/geeorm/log"
	"learn-go/src/projects/geeorm/schema"
	"strings"
)

// Session 数据库会话 负责与数据库交互
type Session struct {
	db       *sql.DB         // 数据库连接
	dialect  dialect.Dialect // 数据库方言
	tx       *sql.Tx         // 事务
	refTable *schema.Schema  // 数据表结构
	clause   clause.Clause   // SQL子句构造
	sql      strings.Builder // SQL语句
	sqlVars  []any           // 占位符对应值
}

// CommonDB DB操作的基本函数集 用于统一管理sql.DB和sql.Tx
type CommonDB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *Session) Raw(sql string, values ...any) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec 封装sql库的基本操作 添加日志打印等
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

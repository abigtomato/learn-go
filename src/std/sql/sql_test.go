package sql

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"testing"
)

func TestSql(t *testing.T) {
	db, _ := sql.Open("sqlite3", "gee.db")
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	_, _ = db.Exec("DROP TABLE IF EXISTS User;")
	_, _ = db.Exec("CREATE TABLE User(Name text);")

	if result, err := db.Exec("INSERT INTO User(`Name`) VALUES (?), (?)", "Tom", "Sam"); err == nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}

	var name string
	if err := db.QueryRow("SELECT Name FROM User LIMIT 1").Scan(&name); err == nil {
		log.Println(name)
	}
}

func TestTx(t *testing.T) {
	db, _ := sql.Open("sqlite3", "gee.db")
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	_, _ = db.Exec("CREATE TABLE IF NOT EXISTS User(`Name` text);")

	tx, _ := db.Begin()
	_, err1 := tx.Exec("INSERT INTO User(`Name`) VALUES (?)", "Tom")
	_, err2 := tx.Exec("INSERT INTO User(`Name`) VALUES (?)", "Jack")
	if err1 != nil || err2 != nil {
		_ = tx.Rollback()
		log.Println("Rollback", err1, err2)
	} else {
		_ = tx.Commit()
		log.Println("Commit")
	}
}

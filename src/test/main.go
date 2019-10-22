package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
)

var DB *sql.DB

func initDb() error {
	var err error
	dsn := "root:root@(localhost:3306)/test2"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(16)
	return nil
}

type User struct {
	Id   int64          `db:"id"`
	Name sql.NullString `db:"string"`
	Email  string       `db:"email"`
}

func testQueryData() {
	sqlstr := "select id, name, email from users where id=?"
	row := DB.QueryRow(sqlstr, 2)
	var user User
	err := row.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%v Email:%d\n", user.Id, user.Name, user.Email)
}


func main(){
	err := initDb()
	if err != nil {
		fmt.Printf("init db failed, err:%v\n", err)
		return
	}
	testQueryData()
}
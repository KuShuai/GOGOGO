package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type user struct {
	id int
	name string
	passward string
	ip string
}

//Go链接MySql
//查询是否有此name
func selectRows()  interface{}{
	sqlStr :=`select count(user) as rows`
	ret,err := db.Exec(sqlStr)
	if err != nil {
		return -1
	}
	return ret
}
func query(param string) bool  {

	//查询单条记录的sql语句
	sqlStr := `select id from user where name =?`
	//执行
	rowObj :=db.QueryRow(sqlStr,param)
	//取结果
	err := rowObj.Scan()
	return err == nil
}

func Insert(name string,passward string,ip string) bool  {
	fmt.Println("55555555")
	//line :=selectRows()
	sqlStr := `insert into user(name,passward,ip) values(?,?,?)`
	ret,err := db.Exec(sqlStr,name,passward,ip)
	if err != nil {
		fmt.Println("1",err)
		return false
	}
	theID,err := ret.LastInsertId()
	if err != nil {
		fmt.Println("2",err)
		return false
	}
	fmt.Println("success",theID)
	return true
}

func OpenSql(err error) error  {

	//数据库信息
	dsn :="root:szkb.123@tcp(192.168.1.13:3306)/zceshi"
	db,err = sql.Open("mysql",dsn)
	if err!=nil {
		return err
	}
	err = db.Ping()
	if err!=nil {
		return err
	}

	return err

}

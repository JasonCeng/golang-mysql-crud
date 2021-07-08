package main

import (
	"golang-mysql-crud/tools"
)

func main() {
	dataSource := "root:123456@tcp(192.168.104.110:33306)/mpc?charset=utf8"
	db,err := tools.InitDb(dataSource)
	tools.CheckErr(err)

	sqlStr := "select id,name,amt from User" //where id = ?
	sqlArgs1 := make([]interface{},0)
	//sqlArgs1 = append(sqlArgs1, "1")
	//rows,err1 := tools.Query(db, sqlStr, sqlArgs1...)
	//tools.CheckErr(err1)


	// 将rows转为Json并打印
	_, _, err = tools.QueryRows2Json(db, sqlStr, sqlArgs1...)
	tools.CheckErr(err)
}

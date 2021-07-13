package main

import (
	"golang-mysql-crud/config"
	"golang-mysql-crud/logging"
	"golang-mysql-crud/rest"
)

func main() {
	// TODO:改造为HTTP服务
	//1.初始化配置
	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	//2.初始化日志
	logging.Initialize()

	//3.初始化Gin框架
	if err := rest.Start(); err != nil {
		panic(err)
	}

	//核心功能：
	////创建数据库连接
	//dataSource := "root:123456@tcp(192.168.104.110:33306)/mpc?charset=utf8"
	//db,err := tools.InitDb(dataSource)
	//tools.CheckErr(err)
	//
	////拼接sql
	//sqlStr := "select id,name,amt from User" //where id = ?
	//sqlArgs1 := make([]interface{},0)
	////sqlArgs1 = append(sqlArgs1, "1")
	////rows,err1 := tools.Query(db, sqlStr, sqlArgs1...)
	////tools.CheckErr(err1)
	//
	////执行查询，并将rows转为json，打印json
	//_, _, err = tools.QueryRows2Json(db, sqlStr, sqlArgs1...)
	//tools.CheckErr(err)
}

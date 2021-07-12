package tools

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func InitDb(dataSource string) (*sql.DB, error) {
	//dataSource = "root:123456@tcp(191.168.104.110:13306)/mpc"
	db, err := sql.Open("mysql",dataSource)
	if err != nil {
		log.Fatalln("sql.Open() error:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln("db.Ping() error:", err)
	}
	return db,err
}

func Query(db *sql.DB, sqlStr string, sqlArgs ...interface{}) (*sql.Rows, error) {
	var (
		id int
		name string
	)
	// sqlStr = "select id, name from User where id = ?"
	rows, err := db.Query(sqlStr, sqlArgs...)
	CheckErr(err)

	for rows.Next() {
		err := rows.Scan(&id, &name)
		CheckErr(err)
		log.Println(id, name)
	}
	err = rows.Err()
	CheckErr(err)

	return rows,err
}

// TODO: 实现 *sql.Rows 到 Json 的转换
// 缺点：返回的值都是string
func QueryRows2Json(db *sql.DB, sqlStr string, sqlArgs ...interface{}) (*sql.Rows, string, error) {
	rows, err := db.Query(sqlStr, sqlArgs...)
	CheckErr(err)
	defer rows.Close()

	columns, err := rows.Columns() // 返回rows的列名
	CheckErr(err)

	count := len(columns)
	tableData := make([]map[string]interface{}, 0) // rows结果集对应的一个map数组,一个map对应一条rows
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ { // 构造一个values数组的地址数组
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{}) // 每一行rows对应的map，key为列名，value为具体的值
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte) // interface{}类型的 val 转为[]byte类型
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry) // 将rows的当前行所对应的一个map插入到tableData中
	}
	jsonData, err := json.Marshal(tableData) // 将map数组转化为json
	CheckErr(err)

	fmt.Println("rows:", rows)
	fmt.Println("rows2json:\n", string(jsonData))
	return rows, string(jsonData), nil
}

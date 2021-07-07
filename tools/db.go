package tools

import (
	"database/sql"
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
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}
	return rows,err
}
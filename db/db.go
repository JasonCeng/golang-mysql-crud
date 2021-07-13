package db

import(
	"database/sql"
	"encoding/base64"
	_  "github.com/go-sql-driver/mysql"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"sync"
    "time"
)

var dbLogger = logging.MustGetLogger("db")

type database struct{
	driverName string
	dataSourceName string
	maxOpenConns int
	maxIdleConns int
	interval time.Duration
}

func initMPCDB(userName string, passwdBase64 string) *database {
	driverName :=viper.GetString("blockchain.database.driverName")
	decodePasswdBytes,_ := base64.StdEncoding.DecodeString(passwdBase64)
	password := string(decodePasswdBytes)
	address :=viper.GetString("blockchain.database.address")
	schema :=viper.GetString("blockchain.database.schema")
	maxOpenConns :=viper.GetInt("blockchain.database.maxOpenConns")
	maxIdleConns :=viper.GetInt("blockchain.database.maxIdleConns")
	interval :=viper.GetDuration("blockchain.database.interval")
	dataSourceName := userName + ":" + password + "@tcp(" + address + ")/" + schema

	db :=&database{
		driverName: driverName,
		dataSourceName: dataSourceName,
		maxOpenConns: maxOpenConns,
		maxIdleConns: maxIdleConns,
		interval: interval,
	}
	go db.checkDBConn()
	return db
}

func (db *database) connectDB() (*sql.DB, error) {
	conn,err :=sql.Open(db.driverName,db.dataSourceName);
	if err!=nil{
		return nil,err
	}

	conn.SetMaxOpenConns(db.maxOpenConns)
	conn.SetMaxIdleConns(db.maxIdleConns)
	return conn,nil
}

var dbOnce sync.Once
var dbConn *sql.DB

func GetDBConn(userName string, passwdBase64 string) (*sql.DB, error) {
	var err error
	dbOnce.Do(func() {
		dbConn,err = initMPCDB(userName, passwdBase64).connectDB()
		if err != nil{
			dbLogger.Errorf("getDBConn err: %s",err)
		}
	})

	return dbConn,err
}

func (db *database) checkDBConn()  {
	dbLogger.Infof("checkDBConn running,interval %s", db.interval)

	for{
		if dbConn ==nil{
			dbLogger.Errorf("db connection is nil")
			var err error
			dbConn,err = db.connectDB()
			if err != nil {
				dbLogger.Errorf("connect db err: %s",err)
				continue
			}
		}

		if err := dbConn.Ping();err != nil{
			dbLogger.Errorf("ping db err: %s",err)
			dbConn.Close()
			var err1 error
			dbConn,err1 =db.connectDB()
			if err1 !=nil{
				dbLogger.Errorf("connect dbConn err: %s",err)
				continue
			}
			dbLogger.Infof("reconnect dbConn successful")
		}
		time.Sleep(db.interval)
	}

}
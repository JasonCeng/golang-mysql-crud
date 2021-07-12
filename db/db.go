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

func initORGSDB() *database {
	driverName :=viper.GetString("blockchain.orgsdatabase.driverName")
	user :=viper.GetString("blockchain.orgsdatabase.user")
	password :=viper.GetString("blockchain.orgsdatabase.password")
	bytes,_ := base64.StdEncoding.DecodeString(password)
	password =string(bytes)
	address :=viper.GetString("blockchain.orgsdatabase.address")
	schema :=viper.GetString("blockchain.orgsdatabase.schema")
	maxOpenConns :=viper.GetInt("blockchain.orgsdatabase.maxOpenConns")
	maxIdleConns :=viper.GetInt("blockchain.orgsdatabase.maxIdleConns")
	interval :=viper.GetDuration("blockchain.orgsdatabase.interval")
	dataSourceName :=user+":"+password+"@tcp("+address+")/"+schema


	db :=&database{
		driverName: driverName,
		dataSourceName: dataSourceName,
		maxOpenConns: maxOpenConns,
		maxIdleConns: maxIdleConns,
		interval: interval,
	}
	go db.checkORGDBConn()
	return db
}

func initMPCDB() *database {
	driverName :=viper.GetString("blockchain.mpcdatabase.driverName")
	user :=viper.GetString("blockchain.mpcdatabase.user")
	password :=viper.GetString("blockchain.mpcdatabase.password")
	bytes,_ := base64.StdEncoding.DecodeString(password)
	password =string(bytes)
	address :=viper.GetString("blockchain.mpcdatabase.address")
	schema :=viper.GetString("blockchain.mpcdatabase.schema")
	maxOpenConns :=viper.GetInt("blockchain.mpcdatabase.maxOpenConns")
	maxIdleConns :=viper.GetInt("blockchain.mpcdatabase.maxIdleConns")
	interval :=viper.GetDuration("blockchain.mpcdatabase.interval")
	dataSourceName :=user+":"+password+"@tcp("+address+")/"+schema


	db :=&database{
		driverName: driverName,
		dataSourceName: dataSourceName,
		maxOpenConns: maxOpenConns,
		maxIdleConns: maxIdleConns,
		interval: interval,
	}
	go db.checkMPCDBConn()
	return db
}

func (db *database)connectDB () (*sql.DB,error) {
	conn,err :=sql.Open(db.driverName,db.dataSourceName);
	if err!=nil{
		return nil,err
	}

	conn.SetMaxOpenConns(db.maxOpenConns)
	conn.SetMaxIdleConns(db.maxIdleConns)
	return conn,nil
}

var mpcOnce sync.Once
var orgOnce sync.Once
var mpcConn *sql.DB
var orgConn *sql.DB

func GetMPCDBConn()*sql.DB{
	mpcOnce.Do(func() {
		var err error
		mpcConn,err = initMPCDB().connectDB()
		if err != nil{
			dbLogger.Errorf("getDBConn err: %s",err)
		}
	})

	return mpcConn
}

func GetORGSDBConn()*sql.DB{
	orgOnce.Do(func() {
		var err error
		orgConn,err = initORGSDB().connectDB()
		if err != nil{
			dbLogger.Errorf("getDBConn err: %s",err)
		}
	})

	return orgConn
}

func (db *database)checkORGDBConn()  {
	dbLogger.Infof("chechORGDBConn running,interval %s",db.interval)

	for{

		if orgConn ==nil{
			dbLogger.Errorf("db connection is nil")
			var err error
			orgConn,err =db.connectDB()
			if err != nil {
				dbLogger.Errorf("connect db err: %s",err)
				continue
			}
		}

		if err :=orgConn.Ping();err!=nil{
			dbLogger.Errorf("ping db err: %s",err)
			orgConn.Close()
			var err1 error
			orgConn,err1 =db.connectDB()
			if err1 !=nil{
				dbLogger.Errorf("connect db err: %s",err)
				continue
			}
			dbLogger.Infof("reconnect db successful")
		}
		time.Sleep(db.interval)
	}

}

func (db *database)checkMPCDBConn()  {
	dbLogger.Infof("chechMPCDBConn running,interval %s",db.interval)

	for{

		if mpcConn ==nil{
			dbLogger.Errorf("db connection is nil")
			var err error
			mpcConn,err = db.connectDB()
			if err != nil {
				dbLogger.Errorf("connect db err: %s",err)
				continue
			}
		}

		if err := mpcConn.Ping();err != nil{
			dbLogger.Errorf("ping db err: %s",err)
			mpcConn.Close()
			var err1 error
			mpcConn,err1 =db.connectDB()
			if err1 !=nil{
				dbLogger.Errorf("connect mpcCon err: %s",err)
				continue
			}
			dbLogger.Infof("reconnect mpcCon successful")
		}
		time.Sleep(db.interval)
	}

}
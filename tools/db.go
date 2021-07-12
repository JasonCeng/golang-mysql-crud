package tools

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"time"
)

const (
	SUCCESS_STATUS = "SUCCESS"
	FAILED_STATUS = "FAILED"
)

type Db struct {
	db *sql.DB
}

type SynReq struct {
	//调用数据同步函数
	SqlStr string `json:"sqlstr,omitempty"`
	TaskInstanceId string `json:"taskInstanceId,omitempty"`
	OrgName string `json:"orgName,omitempty"`
}

type Conf struct {
	ShanghuDatabase ShanghuDatabase
	MpcDatabase MpcDatabase
	Application Application
	Mpcserver Mpcserver
}

type Mpcserver struct {
	Host string
	Port string
}

type Application struct {
	Field1 string
	Field2 string
	Resultfield1 string
	Resultfield2 string
	Resultfield3 string
}

type ShanghuDatabase struct {
	Dbtype string
	Dbname string
	Table1 string
	Table2 string
	Username string
	Password string
	Host string
	Port string
	Ip string
}

type MpcDatabase struct {
	Dbtype string
	Dbname string
	Username string
	Password string
	Host string
	Port string
}

func QuerySyncingData(curdb *sql.DB, ret *Ret, mpcDataSyncStatus *MpcDataSyncStatus) (result *Ret) {
	format := time.Now().Format("2006-01-02")
	startTime := format + " 00:00:00.000000"
	endTime := format + "23:59:59.999999"
	var sql = "select count(*) from mpc.mpc_data_sync_history where tag_id = '" + mpcDataSyncStatus.TagId + "' and status = 'SYNCING' and update_timestmp BETWEEN '" + startTime + "' and '" + endTime + "';"
	var count int64
	countErr := curdb.QueryRow(sql).Scan(&count)
	if countErr != nil {
		onlineLogger.Info(countErr.Error())
	}
	if count > 0 {
		ret.Code = http.StatusBadRequest
		ret.Msg = "DATA SYNCING"
	}
	return ret
}

func SyncTDayData(shanghuConnect, mpcConnect *sql.DB, cisno string, mpcDataSyncStatus *MpcDataSyncStatus) (result *Result, mpcDataSync *MpcDataSyncStatus) {
	result = new(Result)
	dbname := viper.GetString("blockchain.orgsdatabase.schema")
	mpcDb := viper.GetString("blockchain.mpcdatabase.schema")
	mpcTable := viper.GetString("blockchain.mpcdatabase.table")
	mpcCisno := "cisno"
	mpcAmt := "amt"
	mpcDate := "date"

	table1 := viper.GetString("blockchain.orgsdatabase.table1")
	table2 := viper.GetString("blockchain.orgsdatabase.table2")
	field1 := viper.GetString("blockchain.orgsdatabase.field1")
	field2 := viper.GetString("blockchain.orgsdatabase.field2")
	resultfield1 := viper.GetString("blockchain.orgsdatabase.kebianname1")
	resultfield2 := viper.GetString("blockchain.orgsdatabase.countname1")
	resultfield3 := viper.GetString("blockchain.orgsdatabase.timename1")

	unit := viper.GetString("blockchain.orgsdatabase.currency.unit")

	now := time.Now()
	date := now.Format("20060102")

	successNum := 0
	failNum := 0

	result,timeStmp := markSyncing(mpcConnect, mpcDataSyncStatus)

	if table1 != "" && table2 =="" {
		var querySqlFmt = fmt.Sprintf("SELECT %s,%s,%s FROM %s.%s WHERE %s='%s' and %s=%s;",
			resultfield1,
			resultfield2,
			resultfield3,
			dbname,
			table1,
			resultfield1,
			cisno,
			resultfield3,
			date,
		)
		onlineLogger.Info("单表查询sql：", querySqlFmt)
		rows, err := shanghuConnect.Query(querySqlFmt)
		if err != nil {
			onlineLogger.Info("查询数据出现了问题：", err.Error())
			result.Code = http.StatusBadRequest
			result.Msg = err.Error()
			result = markSuccess(mpcConnect, timeStmp, FAILED_STATUS)
			defer rows.Close()
			return result,mpcDataSync
		}
		var deleteSqlFmt = fmt.Sprintf("DELETE FROM %s.%s WHERE %s='%s' and %s=%s;",
			mpcDb,
			mpcTable,
			mpcCisno,
			cisno,
			mpcDate,
			date,
		)
		onlineLogger.Info("单表删除T日该cisno的sql：", deleteSqlFmt)
		result = dropTDayData(mpcConnect, deleteSqlFmt)

		for rows.Next() {
			err := rows.Scan(&resultfield1, &resultfield2, &resultfield3)
			if err != nil {
				onlineLogger.Info("rows Scan failed, detail is ", err.Error())
				result = markSuccess(mpcConnect, timeStmp, FAILED_STATUS)
				defer rows.Close()
				result.Code = http.StatusBadRequest
				result.Msg = err.Error()
				return result,mpcDataSync
			}

			//金额单位转换
			if unit == "1" {
				fAmt, _ := strconv.ParseFloat(resultfield2, 64)
				resultfield2 = strconv.FormatFloat(fAmt/float64(10), 'f', -1, 64)
			} else if unit == "2" {
				fAmt, _ := strconv.ParseFloat(resultfield2, 64)
				resultfield2 = strconv.FormatFloat(fAmt/float64(100), 'f', -1, 64)
			}
			//插入数据
			var insertSqlFmt = fmt.Sprintf("INSERT INTO %s.%s(%s,%s,%s) values(?,?,?);",
				mpcDb,
				mpcTable,
				mpcCisno,
				mpcAmt,
				mpcDate,
			)
			_, err = mpcConnect.Exec(insertSqlFmt, resultfield1, resultfield2, resultfield3)
			//onlineLogger.Info("插入sql:", insertSqlFmt)
			if err != nil {
				failNum++
			} else {
				successNum++
			}
		}

		result.Fail = "数据同步失败" + strconv.Itoa(failNum) + "条"
		result.Success = "成功插入" + strconv.Itoa(successNum) + "条"

		if failNum != 0 {
			mpcDataSyncStatus.Status = FAILED_STATUS
			result = markSuccess(mpcConnect, timeStmp, FAILED_STATUS)
		} else {
			mpcDataSyncStatus.Status = SUCCESS_STATUS
			result = markSuccess(mpcConnect, timeStmp, SUCCESS_STATUS)
		}

		defer rows.Close()
	}

	if table1 != "" && table2 != "" {
		//1.取数
		var querySqlFmt = fmt.Sprintf("SELECT %s,%s,%s FROM %s t1, %s t2 WHERE t1.%s=t2.%s and %s='%s' and %s=%s;",
			resultfield1,
			resultfield2,
			resultfield3,
			table1,
			table2,
			field1,
			field2,
			resultfield1,
			cisno,
			resultfield3,
			date,
		)
		onlineLogger.Info("单表查询sql：", querySqlFmt)
		rows, err := shanghuConnect.Query(querySqlFmt)
		if err != nil {
			onlineLogger.Info("多表查询数据出现了问题：", err.Error())
			result.Fail = "查询数据出现了问题，fail"
			result.Code = http.StatusBadRequest
			result.Msg = err.Error()
			result = markSuccess(mpcConnect, timeStmp, FAILED_STATUS)
			defer rows.Close()
			return result,mpcDataSync
		}
		var deleteSqlFmt = fmt.Sprintf("DELETE FROM %s.%s WHERE %s='%s' and %s=%s;",
			mpcDb,
			mpcTable,
			mpcCisno,
			cisno,
			mpcDate,
			date,
		)
		onlineLogger.Info("单表删除T日该cisno的sql：", deleteSqlFmt)
		result = dropTDayData(mpcConnect, deleteSqlFmt)

		for rows.Next() {
			err := rows.Scan(&resultfield1, &resultfield2, &resultfield3)
			if err != nil {
				onlineLogger.Info("rows Scan failed, detail is ", err.Error())
				result = markSuccess(mpcConnect, timeStmp, FAILED_STATUS)
				defer rows.Close()
				result.Code = http.StatusBadRequest
				result.Msg = err.Error()
				return result,mpcDataSync
			}

			//金额单位转换
			if unit == "1" {
				fAmt, _ := strconv.ParseFloat(resultfield2, 64)
				resultfield2 = strconv.FormatFloat(fAmt/float64(10), 'f', -1, 64)
			} else if unit == "2" {
				fAmt, _ := strconv.ParseFloat(resultfield2, 64)
				resultfield2 = strconv.FormatFloat(fAmt/float64(100), 'f', -1, 64)
			}
			//插入数据
			var insertSqlFmt = fmt.Sprintf("INSERT INTO %s.%s(%s,%s,%s) values(?,?,?);",
				mpcDb,
				mpcTable,
				mpcCisno,
				mpcAmt,
				mpcDate,
			)
			_, err = mpcConnect.Exec(insertSqlFmt, resultfield1, resultfield2, resultfield3)
			if err != nil {
				failNum++
			} else {
				successNum++
			}
		}

		result.Fail = "数据同步失败" + strconv.Itoa(failNum) + "条"
		result.Success = "成功插入" + strconv.Itoa(successNum) + "条"

		if failNum != 0 {
			mpcDataSyncStatus.Status = FAILED_STATUS
			result = markSuccess(mpcConnect, timeStmp, FAILED_STATUS)
		} else {
			mpcDataSyncStatus.Status = SUCCESS_STATUS
			result = markSuccess(mpcConnect, timeStmp, SUCCESS_STATUS)
		}

		defer rows.Close()
	}

	mpcDataSyncStatus.TagId = date + cisno

	var counts string
	var countSqlFmt = fmt.Sprintf("SELECT COUNT(1) FROM %s.%s WHERE %s='%s' AND %s=%s;",
		mpcDb,
		mpcTable,
		mpcCisno,
		cisno,
		mpcDate,
		date,
	)
	onlineLogger.Info("计算同步数量的sql：", countSqlFmt)
	result, counts = queryCount(mpcConnect, countSqlFmt)
	mpcDataSyncStatus.Counts = counts

	nowTimeStmp := time.Now().Format("2006-01-02 15:04:05.000000")
	mpcDataSyncStatus.InsertTimestmp = nowTimeStmp

	result = markSuccess(mpcConnect, timeStmp, "SUCCESS")

	result.Code = http.StatusOK
	return result,mpcDataSyncStatus
}

func Get(curdb *sql.DB, sqlstr string, ret *Ret) (result *Ret, sumAmt sql.NullFloat64) {
	var mpcData *mpcDataInfo = new(mpcDataInfo)
	rows := curdb.QueryRow(sqlstr)

	err := rows.Scan(&mpcData.SumAmt)
	if err != nil {
		ret.Code = http.StatusBadRequest
		ret.Msg = err.Error()
		return ret,sumAmt
	}
	return ret,mpcData.SumAmt
}

func markSyncing(mpcConnect *sql.DB, mpcDataSyncStatus *MpcDataSyncStatus) (result *Result, timeStmp string) {
	result = new(Result)
	timeStmp = time.Now().Format("2006-01-02 15:04:05.000000")
	_, err := mpcConnect.Exec("insert into mpc.mpc_data_sync_history(tag_id,status,update_timestmp) values (?,?,?)", mpcDataSyncStatus.TagId,"SYNCING",timeStmp)
	if err != nil {
		onlineLogger.Info(err.Error())
		result.Code = http.StatusBadRequest
		result.Msg = err.Error()
		return result,timeStmp
	}
	result.Code = http.StatusOK
	result.Msg = "记录数据同步状态为SYNCING成功"
	return result,timeStmp
}

func dropTDayData(mpcConnect *sql.DB, deleteSql string) (result *Result) {
	result = new(Result)
	_, e := mpcConnect.Exec(deleteSql)
	if e != nil {
		onlineLogger.Info("删除当天数据出现了问题：", e.Error())
		result.Code = http.StatusBadRequest
		result.Msg = e.Error()
		return result
	}
	result.Code = http.StatusOK
	result.Msg = "删除当天数据成功"
	return result
}

func markSuccess(mpcConnect *sql.DB, timeStmp, status string) (result *Result) {
	result = new(Result)
	updateSql := "update mpc.mpc_data_sync_history set status = ? where update_timestmp = ?"
	_, err := mpcConnect.Exec(updateSql, status, timeStmp)
	if err != nil {
		onlineLogger.Info(err.Error())
		result.Code = http.StatusBadRequest
		result.Msg = err.Error()
		return result
	}
	result.Code = http.StatusOK
	result.Msg = "记录数据同步状态为" + status + "成功"
	return result
}

func queryCount(mpcConnect *sql.DB, countSql string) (result *Result, counts string) {
	result = new(Result)
	var count int64
	countErr := mpcConnect.QueryRow(countSql).Scan(&count)
	if countErr != nil {
		onlineLogger.Info("查询总条数出现了问题：", countErr.Error())
		result.Code = http.StatusBadRequest
		result.Msg = countErr.Error()
		return result,"null"
	}
	counts = strconv.FormatInt(count,10)
	result.Code = http.StatusOK
	result.Msg = "计算同步的数据总条数成功"
	return result,counts
}
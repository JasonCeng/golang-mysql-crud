package tools

import (
	"database/sql"
	"encoding/json"
	"github.com/op/go-logging"
	"golang-mysql-crud/db"
	"net/http"
	"strconv"
	"time"
	//"mpc-data-tools/mpc-online-data/db"
)

var onlineLogger = logging.MustGetLogger("app")

type Ret struct {
	Code int
	Msg string
}

type Result struct {
	Code int
	Msg string
	Success string
	Fail string
}

type MpcDataSyncStatus struct {
	TaskInstanceId string
	OrgName string
	Counts string
	Hash string
	Status string
	TagId string
	InsertTimestmp string
}

const (
	METHOD_GET = "GET"
	METHOD_POST = "POST"
)

const (
	JSON_TYPE = "application/json; charset=utf-8"
)

type ChaincodeReq struct {
	Channel string
	ChaincodeName string
	FunctionName string
	Args *reqArgs
}

type reqArgs struct {
	OrgName string
	TaskInstanceId string
	Counts string
	DataHash string
	Status string
	TagId string
	SyncType string
}

type mpcDataInfo struct {
	Cisno  string          `db:"cisno"`
	Amt    int             `db:"amt"`
	Date   string          `db:"date"`
	SumAmt sql.NullFloat64 `db:"sum(amt)"`
}

func MpcData(synReq *SynReq) (ret *Ret) {
	onlineLogger.Info("============ mpc-online-data start ============")
	startTime := time.Now()

	taskInstanceId := synReq.TaskInstanceId

	orgName := synReq.OrgName
	sqlstr := synReq.SqlStr
	cisno := getCisno(sqlstr)
	onlineLogger.Info("taskInstanceId=", taskInstanceId, " orgName=", orgName, " cisno=", cisno, " sqlstr=", sqlstr)
	readReqEndTime := time.Now()
	onlineLogger.Info("读取参数耗时：", readReqEndTime.Sub(startTime))

	mpcDataSyncStatus := new(MpcDataSyncStatus)
	mpcDataSyncStatus.TaskInstanceId = taskInstanceId
	now := time.Now()
	date := now.Format("20060102")
	mpcDataSyncStatus.TagId = date + cisno
	mpcDataSyncStatus.OrgName = orgName

	shanghuConnect := db.GetORGSDBConn()
	mpcConnect := db.GetMPCDBConn()
	connectDBTime := time.Now()
	onlineLogger.Info("连接数据库耗时：", connectDBTime.Sub(readReqEndTime))

	ret = new(Ret)
	ret.Code = http.StatusOK

	ret = QuerySyncingData(mpcConnect, ret, mpcDataSyncStatus)
	isSyncingTime := time.Now()
	onlineLogger.Info("判断是否有SYNCING状态的数据耗时：", isSyncingTime.Sub(connectDBTime))

	if ret.Code == http.StatusOK {
		onlineLogger.Info("============= sync data start ==============")
		result, mpcData := SyncTDayData(shanghuConnect, mpcConnect, cisno, mpcDataSyncStatus)
		queryDataTime := time.Now()
		onlineLogger.Info("同步数据耗时：", queryDataTime.Sub(isSyncingTime))
		onlineLogger.Info("============= sync data end ==============")
		onlineLogger.Info(result.Success, result.Fail)

		onlineLogger.Info("============= mpc get data start ==============")
		ret, amtResult := Get(mpcConnect, sqlstr, ret)
		if ret.Code == http.StatusBadRequest{
			ret.Msg = "0"
		}else {
			ret.Msg = strconv.FormatFloat(amtResult.Float64, 'f', -1, 64)
		}

		ret_json, _ := json.Marshal(ret)
		onlineLogger.Info("ret_json:", ret_json)
		onlineLogger.Info("============ mpc get data start ==============")

		go func() {
			var hashRes string
			result, hashRes = hashSyncData(mpcConnect, cisno)
			mpcDataSyncStatus.Hash = hashRes

			chainCode(result, mpcData, ret)
			chainCodeTime := time.Now()
			onlineLogger.Info("上链耗时：", chainCodeTime.Sub(queryDataTime))
		}()
	}

	onlineLogger.Info("判断是否有SYNCING状态的数据ret：%v", ret)
	onlineLogger.Info("============ mpc-online-data end ==============")

	endTime := time.Now()
	spendTime := endTime.Sub(startTime)
	onlineLogger.Info("总耗时：", spendTime)
	return ret
}
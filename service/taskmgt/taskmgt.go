package taskmgt

import (
	"fmt"
	"github.com/op/go-logging"
	"golang-mysql-crud/db"
	"golang-mysql-crud/tools"
	"strings"
)

const (
	STATUS_OK = "OK"
	STATUS_ERROR = "ERROR"
)

var taskMgtLogger = logging.MustGetLogger("taskMgt")

type taskMgt struct {
	txId string
}

func GetTaskMgt(txId string) *taskMgt {
	return &taskMgt{txId: txId}
}

func (task *taskMgt) TaskService(queryReq *tools.QueryReq) *tools.QueryResult {
	result := task.QueryData(queryReq)
	return result
}

func (task *taskMgt) QueryData(queryReq *tools.QueryReq) *tools.QueryResult {
	//1、接收请求参数
	userName := queryReq.UserName
	passwdBase64 := queryReq.PasswdBase64
	tableName := queryReq.TableName
	queryFields := queryReq.QueryFields
	queryCondition := queryReq.QueryCondition

	//2、获取数据库连接
	dbConn,err := db.GetDBConn(userName, passwdBase64)
	tools.CheckErr(err)

	//3、拼接sql
	queryFieldsStr := strings.Join(queryFields, ",")
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s;", queryFieldsStr, tableName, queryCondition)
	//sqlArgs1 := make([]interface{},0)

	//4、执行查询并返回json字符串
	_, resultJsonStr, err := tools.QueryRows2Json(dbConn, sqlStr)
	tools.CheckErr(err)

	return &tools.QueryResult {
		Status: STATUS_OK,
		Message: "查询数据成功",
		Uuid: tools.GenerateUUID(),
		QueryJsonData: resultJsonStr,
	}
}
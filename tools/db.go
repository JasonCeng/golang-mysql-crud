package tools

import (
	"database/sql"
)

const (
	SUCCESS_STATUS = "SUCCESS"
	FAILED_STATUS = "FAILED"
)

const (
	SUCCESS_CODE    = 0
	CONNECT_DB_FAIL = -11130
)

type Db struct {
	db *sql.DB
}

type QueryReq struct {
	UserName string `json:"userName,omitempty"`
	PasswdBase64 string `json:"passwdBase64,omitempty"`
	TableName string `json:"tableName,omitempty"`
	QueryFields []string `json:"queryFields,omitempty"` //动态属性：查询字段
	QueryCondition string `json:"queryCondition,omitempty"` //查询条件
}

type QueryResult struct {
	Status string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
	Uuid string `json:"uuid,omitempty"`
	QueryJsonData string `json:"queryJsonData,omitempty"`
}

type Error struct {
	Code int `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
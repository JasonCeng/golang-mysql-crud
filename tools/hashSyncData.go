package tools

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"time"
)

func Hash256(srcs ...string)string  {
	hashFunc :=sha256.New()
	for _,src:= range srcs{
		hashFunc.Write([]byte(src))
	}
	return  hex.EncodeToString(hashFunc.Sum(nil))
}

func hashSyncData(mpcConnect *sql.DB,cisno string)(result *Result,hashRes string){
	result = new(Result)
	now :=time.Now()
	date :=now.Format("20060102")
	var hashStr string
	var querySql ="select cisno,amt,date from mpc.customer where cisno='"+ cisno +"' and date="+ date +" order by amt,cisno,date"+";"
	onlineLogger.Info("查詢T日的該聯合計算的cino 全部数据sql",querySql)
	rows,err :=mpcConnect.Query(querySql)
	if err != nil{
		onlineLogger.Info("查询数据出现问题",err.Error())
		result.Code =http.StatusBadRequest
		result.Msg =err.Error()
		defer rows.Close()
		return  result,"null"
	}

	resultfield1 := ""
	resultfield2 := ""
	resultfield3 := ""
	for rows.Next(){
		err :=rows.Scan(&resultfield1,&resultfield2,&resultfield2)
		if err!=nil{
			onlineLogger.Info("rows scan failed: detail is",err.Error())
			result.Code =http.StatusBadRequest
			result.Msg =err.Error()
			defer rows.Close()
			return  result,"null"
		}
		hashStr = hashStr + resultfield1+resultfield2+resultfield3
	}
	hash256 :=Hash256(hashStr)
	hashRes = hash256
	defer rows.Close()
	result.Code = http.StatusOK
	return result,hashRes
}

package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-mysql-crud/tools"
	"net/http"
)

const (
	REST_VERSION = "1.0"
	ERR_DEFAULT = "SYSTEM ERROR"
)

const (
	SUCCESS_CPDE    = 0
	SYSTEM_ERR_CODE = -10000
	URL_NOT_FOUND   = -10001
	USER_LOGIN_FAIL = -20000
	USER_NOT_LOGIN  = -20001
	TASK_RUN_STATUS = -30000
	START_TASK_FAIL = -30001
	EXIST_STATUS    = -40000
)

type response struct {
	TxId string       `json:"txId,omitempty"`
	Version string    `json:"version,omitempty"`
	Code int          `json:"code"`
	Msg string        `json:"msg,omitempty"`
}

func noRouterHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		warnMsg := fmt.Sprintf("%s not found.", context.Request.URL.Path)
		warnHandler(context, http.StatusNotFound, URL_NOT_FOUND, warnMsg)
	}
}

func errHandler(context *gin.Context, httpCode, errorCode int, errMsg string)  {
	resp := response{
		TxId:    context.GetString("txId"),
		Version: REST_VERSION,
		Code:    errorCode,
		Msg:     errMsg,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		restLogger.Error(err.Error())
		return
	}

	if data != nil {
		context.Data(httpCode, tools.JSON_TYPE, data)
	}else {
		context.Data(httpCode, tools.JSON_TYPE, []byte(ERR_DEFAULT))
	}

	context.Set("resp", errMsg)
	restLogger.Error(errMsg)
}

func warnHandler(context *gin.Context, httpCode, errorCode int, warnMsg string)  {
	resp := response{
		TxId:    context.GetString("txId"),
		Version: REST_VERSION,
		Code:    errorCode,
		Msg:     warnMsg,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		errHandler(context, http.StatusBadRequest, SYSTEM_ERR_CODE, err.Error())
		return
	}

	context.Data(httpCode, tools.JSON_TYPE, data)
	context.Set("resp", warnMsg)
	restLogger.Warning(warnMsg)
}

func okHandler(context *gin.Context, okMsg string)  {
	resp := response{
		TxId:    context.GetString("txId"),
		Version: REST_VERSION,
		Code:    SUCCESS_CPDE,
		Msg:     okMsg,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		errHandler(context, http.StatusBadRequest, SYSTEM_ERR_CODE, err.Error())
		return
	}

	context.Data(http.StatusOK, tools.JSON_TYPE, data)
	context.Set("resp", okMsg)
}
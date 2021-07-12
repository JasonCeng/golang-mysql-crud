package rest

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"io/ioutil"
	"bytes"
	"time"
	"encoding/json"
	//"mpc-data-tools/mpc-online-data/service/sessionmgt"
	"golang-mysql-crud/tools"
)

type RequestLog struct {
	TxId string
	Url string
	Method string
	ReqBody string
}

type ResponseLog struct {
	TxId string
	Resp string
	HandleTime string
}

func restLog() gin.HandlerFunc {
	return func(context *gin.Context) {
		reqBody, _ := ioutil.ReadAll(context.Request.Body)

		startTime := time.Now()
		txId := tools.GenerateUUID()
		context.Set("txId", txId)

		reqLog := RequestLog{
			TxId:    txId,
			Url:     context.Request.URL.String(),
			Method:  context.Request.Method,
			ReqBody: string(reqBody),
		}

		reqLogInfo, err := json.Marshal(reqLog)
		if err != nil {
			restLogger.Errorf("requesting Marshal err: %s", err)
		}

		restLogger.Infof("handle request-> %s", reqLogInfo)

		context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

		context.Next()

		respLog := ResponseLog{
			TxId:       txId,
			Resp:       context.GetString("resp"),
			HandleTime: fmt.Sprintf("%s", time.Since(startTime)),
		}

		respLogInfo, err := json.Marshal(respLog)
		if err != nil {
			restLogger.Errorf("responseLog Marshal err: %s", err)
		}

		restLogger.Infof("handle result-> %s", respLogInfo)
	}
}

func setRespType() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "accept, content-type")
	}
}
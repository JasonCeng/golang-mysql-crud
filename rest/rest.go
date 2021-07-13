package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"golang-mysql-crud/service/taskmgt"
	"golang-mysql-crud/tools"
	"net/http"
)

var restLogger = logging.MustGetLogger("rest")

func Start() error {

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(restLog())
	router.Use(setRespType())

	parent := router.Group("/rest")

	parent.GET("/db", QueryData())

	router.NoRoute(noRouterHandler())

	host := viper.GetString("blockchain.rest.listenaddress")
	restLogger.Infof("blockchain-rest-db starting, listen on %s", host)

	return router.Run(host)
}

func QueryData() gin.HandlerFunc {
	return func(ct *gin.Context) {
		var queryReq tools.QueryReq

		err := ct.BindJSON(&queryReq)
		if err != nil {
			errMsg := "queryReq BindJSON err: " + err.Error()
			errHandler(ct, http.StatusBadRequest, SYSTEM_ERR_CODE, errMsg)
			return
		}

		QueryResult := taskmgt.GetTaskMgt(ct.GetString("txId")).TaskService(&queryReq)
		if QueryResult.Status != "OK" {
			errHandler(ct, http.StatusBadRequest, QueryResult.Error.Code, QueryResult.Error.Message)
			return
		}

		okHandler(ct, QueryResult.Message)
	}
}
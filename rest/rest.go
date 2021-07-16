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
	//发版模式创建路由对象
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	//设置rest日志中间件
	router.Use(restLog())
	//设置跨域访问
	router.Use(setRespType())

	//创建Group
	parent := router.Group("/rest/db")

	//前置服务通用数据查询 POST请求接口
	parent.POST("/queryData", QueryData())

	//URL_NOT_FOUND处理
	router.NoRoute(noRouterHandler())

	//读取配置
	host := viper.GetString("blockchain.rest.database.listenaddress")
	restLogger.Infof("blockchain-rest-db starting, listen on %s", host)

	//无限期启动HTTP服务和监听，除非有报错发生
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

		okHandler(ct, QueryResult)
	}
}
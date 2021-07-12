package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"golang-mysql-crud/tools"
)

var restLogger = logging.MustGetLogger("rest")

func Start() error {

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(restLog())
	router.Use(setRespType())

	parent := router.Group("/mpc/api")

	parent.GET("/getMpcData", MpcData())

	router.NoRoute(noRouterHandler())

	host := viper.GetString("blockchain.rest.listenaddress")
	restLogger.Infof("blockchain-mpc starting, listen on %s", host)

	return router.Run(host)
}

func MpcData() gin.HandlerFunc {
	return func(ct *gin.Context) {
		var synReq tools.SynReq

		synReq.SqlStr = ct.Query("sqlstr")
		synReq.TaskInstanceId = ct.Query("taskInstanceId")
		synReq.OrgName = ct.Query("orgName")

		result := tools.MpcData(&synReq)
		okHandler(ct, result.Code, result.Msg)
	}
}
package tools

import (
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)
var httpLogger = logging.MustGetLogger("tools")
type RpcResponse struct{
	Msg string `json:"msg,omitempty"`
}

type HttpClient struct{
	client * http.Client
}

func GetHttpClient()(*HttpClient,error)  {
	httpClient :=&HttpClient{}
	client := &http.Client{
		Timeout: time.Duration(10*time.Second),
	}
	httpClient.client = client
	return  httpClient,nil
}

func (hc*HttpClient)HttpRequest(method,contentType,url string,msg[]byte) ([]byte,error) {
	req,err :=http.NewRequest(method,url,strings.NewReader(string(msg)))
	if err !=nil{
		return  nil,err
	}
	req.Header.Set("Content-Type",contentType)
	resp,err := hc.client.Do(req)
	if err!=nil{
		return nil,err
	}
	defer resp.Body.Close()
	if resp.StatusCode !=http.StatusOK{
		errMsg,err := ioutil.ReadAll(resp.Body)
		if err !=nil{
			return nil,err
		}
		return nil,fmt.Errorf("status code:%d,errMsg:%s",resp.StatusCode,errMsg)
	}
  data,err := ioutil.ReadAll(resp.Body)
  if err !=nil{
  	return nil,fmt.Errorf("read response err:%s",err)
  }

  return  data,nil
}

func Dorequest(chaincodeReq *ChaincodeReq,ret *Ret)  {

	host :=viper.GetString("blockchain.mpcserve.host")
	port :=viper.GetString("blockchain.mpcserve.port")
	url :=fmt.Sprintf("http://%s:%s/pub?topic=audit",host,port)
	onlineLogger.Info(url)
	sendData,err :=json.Marshal(chaincodeReq)
	if err!=nil {
		onlineLogger.Errorf("json解析失败: %v",err)
		return
	}
	onlineLogger.Infof("存证:%s",string(sendData))
	httpClient, _ := GetHttpClient()
	resp,err := httpClient.HttpRequest(METHOD_POST,JSON_TYPE,url,sendData)
	httpLogger.Infof("httpClient.HttpRequest resp: %s",resp)
	if err!= nil{
		onlineLogger.Errorf("上链失败: %v",err)
	}
	rep :=string(resp)
	if rep !="OK"{
		httpLogger.Errorf("request(%s) err:%s",url,rep)
	}else{
		httpLogger.Infof("send StartTask Msg successful,req:%s",sendData)
	}
	return
}

func chainCode(result *Result,mpcData *MpcDataSyncStatus,ret *Ret)  {

	if result.Code == http.StatusOK{
		chaincodeReq :=&ChaincodeReq{
			Channel :"mpcchannel",
			ChaincodeName :"audit_cc",
			FunctionName :"saveStatusOfDataSync",
			Args:&reqArgs{
				Counts:mpcData.Counts,
				TagId: mpcData.TagId,
				Status: mpcData.Status,
				DataHash : mpcData.Hash,
				OrgName: mpcData.OrgName,
				TaskInstanceId:mpcData.TaskInstanceId,
				SyncType:"2",
			},
		}

		onlineLogger.Info("==========invoke saveStatusOfDataSync start=============")
		Dorequest(chaincodeReq,ret)
		onlineLogger.Info("==========invoke saveStatusOfDataSync start=============")
	}else {
		ret.Code = http.StatusBadRequest
		ret.Msg = result.Msg
	}

}
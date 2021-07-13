package tools

import (
	"fmt"
	"github.com/op/go-logging"
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
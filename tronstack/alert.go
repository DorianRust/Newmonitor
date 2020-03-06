package tronstack

import (
	"github.com/sasaxie/monitor/common/database/influxdb"
	"time"
	"net/http"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"fmt"
	"strings"
)


type TronStackAlert struct {
}

func (alert TronStackAlert) CheckEventQueryStatus(ip string) {
	healthCheck := ip + "/healthcheck"
	status := checkHeathForEventQuery(healthCheck)
	httpMap := map[string]string{
		"IP":   ip,
	}
	httpFields := map[string]interface{}{
		"IP":      ip,
		"Status":   status,
	}
	influxdb.Client.WriteByTime("api_report_eventQuery_status", httpMap, httpFields, time.Now())
}

func (alert TronStackAlert)CheckTronStackStatus(ip string) {
	status := checkHeathForTronStack(ip + "/wallet/getnowblock")
	httpMap := map[string]string{
		"IP":   ip,
	}
	httpFields := map[string]interface{}{
		"IP":      ip,
		"Status":   status,
	}
	influxdb.Client.WriteByTime("api_report_tronStack_status", httpMap, httpFields, time.Now())
}

func (alert TronStackAlert)CheckTronStackResponse(url string, iplist []string) {
	totalResponse2xx, totalResponse4xx, totalResponse5xx:= int64(0), int64(0), int64(0)
	for _, ip := range iplist {
		response2xx, response4xx, response5xx := GetResponseCount(ip + "/monitor")
		httpMap := map[string]string{
			"IP":   ip,
		}
		httpFields := map[string]interface{}{
			"IP":      ip,
			"Response2xx":   response2xx,
			"Response4xx":   response4xx,
			"Response5xx":   response5xx,
		}
		totalResponse2xx += response2xx
		totalResponse4xx += response4xx
		totalResponse5xx += response5xx
		influxdb.Client.WriteByTime("api_report_tronStack_response", httpMap, httpFields, time.Now())
	}
	httpMap := map[string]string{
		"IP":   url,
	}
	httpFields := map[string]interface{}{
		"IP":      url,
		"Response2xx":   totalResponse2xx,
		"Response4xx":   totalResponse4xx,
		"Response5xx":   totalResponse5xx,
	}
	influxdb.Client.WriteByTime("api_report_tronStack_response", httpMap, httpFields, time.Now())
}

func checkHeathForEventQuery(httpUrl string) string {
	result, err := http.Get(httpUrl)

	if err != nil {
		logs.Error("error", err)
		return "Down"
	}
	defer func(result  *http.Response) {
		if (result != nil){
			result.Body.Close()
		}

	}(result)
	s, _ := ioutil.ReadAll(result.Body) //把  body 内容读入字符串 s
	var content string
	content = fmt.Sprintf("%s", s)     //在返回页面中显示内容。
	if (strings.EqualFold(content , "OK")){
		return "OK"
	}
	return  "Down"
}


func checkHeathForTronStack(httpUrl string) string {
	result, err := http.Get(httpUrl)

	if err != nil {
		return "DOWN"
	}
	defer func(result  *http.Response) {
		if (result != nil){
			result.Body.Close()
		}

	}(result)
	s, _ := ioutil.ReadAll(result.Body) //把  body 内容读入字符串 s
	var content string
	content = fmt.Sprintf("%s", s)     //在返回页面中显示内容。
	if (result.StatusCode == 200 && len(content) != 0){
		return "OK"
	}
	return  "Down"
}


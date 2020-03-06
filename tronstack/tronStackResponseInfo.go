package tronstack

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
)

type ApiDetailInfo struct {
	Name string `json:"name"`
	Count int64 `json:"count"`
	FailCount int64 `json:"failCount"`
	Count4xx int64  `json:"count4xx"`
	Count5xx int64  `json:"count5xx"`
	Count2xx int64  `json:"count2xx"`
}

type ApiInfo struct {
	TotalCount int64 `json:"totalCount"`
	TotalFailCount int64 `json:"totalFailCount"`
	TotalCount2xx int64 `json:"totalCount2xx"`
	TotalCount4xx int64 `json:"totalCount4xx"`
	TotalCount5xx int64 `json:"totalCount5xx"`
	Detail  []ApiDetailInfo `json:"apiDetailInfo"`
}

type  DisconnectionDetailInfo struct {
	Reason string `json:"reason"`
	Count int64 `json:"count"`
}

type LatencyInfo struct {
	Top99 int64 `json:"top99"`
	Top95 int64 `json:"top95"`
	TotalCount int64 `json:"totalCount"`
	Delay1S int64 `json:"delay1S"`
	Delay2S int64 `json:"delay2S"`
	Delay3S int64 `json:"delay3S"`
	Detail LatencyDetailInfo `json:"detail"`
}


type LatencyDetailInfo struct {
	Witness string `json:"witness"`
	Top99 int64 `json:"top99"`
	Top95 int64 `json:"top95"`
	Count int64 `json:"count"`
	Delay1S int64 `json:"delay1S"`
	Delay2S int64 `json:"delay2S"`
	Delay3S int64 `json:"delay3S"`
}

type NetInfo struct {
	ErrorProtoCount int64 `json:"errorProtoCount"`
	Api ApiInfo `json:"api"`
	ConnectionCount int64 `json:"connectionCount"`
	ValidConnectionCount int64 `json:"validConnectionCount"`
	TCPInTraffic int64 `json:"TCPInTraffic"`
	TCPOutTraffic int64 `json:"TCPOutTraffic"`
	DisconnectionCount int64 `json:"disconnectionCount"`
	DisconnectionDetail []*DisconnectionDetailInfo `json:"disconnectionDetail"`
	UDPInTraffic int64 `json:"UDPInTraffic"`
	UDPOutTraffic int64 `json:"UDPOutTraffic"`
	Latency LatencyInfo `json:"latency"`
}

type DataInfo struct {
	Interval int64 `json:"interval"`
	NetInfo NetInfo `json:"netInfo"`
}

type TronStackResponseInfo struct {
	Status int64 `json:"status"`
	Msg string `json:"msg"`
	Data DataInfo `json:"dataInfo"`
}

func GetResponseCount(httpUrl string) (int64, int64, int64) {
	result, err := http.Get(httpUrl)
	bytes, err := ioutil.ReadAll(result.Body)
	fmt.Println(string(bytes))

	if err != nil {
		fmt.Print("error", err)
		return 0, 0, 0
	}
	defer func(result  *http.Response) {
		if (result != nil){
			result.Body.Close()
		}

	}(result)

	var res  =  &TronStackResponseInfo{}
	json.Unmarshal(bytes, &res)
	return res.Data.NetInfo.Api.TotalCount2xx, res.Data.NetInfo.Api.TotalCount4xx, res.Data.NetInfo.Api.TotalCount5xx
}

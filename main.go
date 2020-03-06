package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron"
	"github.com/sasaxie/monitor/alerts"
	"github.com/sasaxie/monitor/common/config"
	"github.com/sasaxie/monitor/common/database/influxdb"
	"github.com/sasaxie/monitor/datamanger"
	"github.com/sasaxie/monitor/reports"
	_ "github.com/sasaxie/monitor/routers"
	"time"
	"github.com/sasaxie/monitor/slack"
	"github.com/sasaxie/monitor/tronstack"
)

var httpex = "https://httpapi.tronex.io/"

var httpevent = "https://api.tronex.io/"

var httpTestEvent = "https://testapi.tronex.io/"

var fullUrl =
	[]string{
	          "wallet/getnowblock",
	          "wallet/getaccount?address=TQJkDDDGQoi2yrfdpG5nUSHcgJ1KpBXan7&visible=true",
	          "wallet/listnodes",
	          "wallet/getblockbynum?num=1",
	          "wallet/getblockbyid?value=00000000000000010ff5414c5cfbe9eae982e8cef7eb2399a39118e1206c8247",
		      "wallet/getblockbylimitnext?startNum=1&endNum=1",
		      "wallet/getblockbylatestnum?num=1",
		      "wallet/gettransactionbyid?value=6a9c700b66caf4a47ad6b1c3cc06dbcf04d1ce98ad2f305bba23918cd06fcee9",
		      "wallet/gettransactioninfobyid?value=6a9c700b66caf4a47ad6b1c3cc06dbcf04d1ce98ad2f305bba23918cd06fcee9",
		      "wallet/gettransactioncountbyblocknum?num=4000012",
		      "wallet/listwitnesses",
		      "wallet/getassetissuelist",
		      "wallet/getnextmaintenancetime",
		      "wallet/listproposals",
		      "wallet/getproposalbyid?id=1",
		      "wallet/getexchangebyid?id=1",
		      "wallet/listexchanges",
		      "wallet/getchainparameters",
		      "wallet/getnodeinfo",
	        }
var solidityUrl =
	[]string{
	           "walletsolidity/getnowblock",
	           "walletsolidity/getaccount?address=TQJkDDDGQoi2yrfdpG5nUSHcgJ1KpBXan7&visible=true",
		       "walletsolidity/getblockbynum?num=1",
		       "walletsolidity/getblockbyid?value=00000000000000010ff5414c5cfbe9eae982e8cef7eb2399a39118e1206c8247",
		       "walletsolidity/getblockbylimitnext?startNum=1&endNum=1",
		       "walletsolidity/getblockbylatestnum?num=1",
		       "walletsolidity/gettransactionbyid?value=6a9c700b66caf4a47ad6b1c3cc06dbcf04d1ce98ad2f305bba23918cd06fcee9",
		       "walletsolidity/gettransactioninfobyid?value=6a9c700b66caf4a47ad6b1c3cc06dbcf04d1ce98ad2f305bba23918cd06fcee9",
		       "walletsolidity/gettransactioncountbyblocknum?num=4000012",
		       "walletsolidity/listwitnesses",
		       "walletsolidity/getassetissuelist",
		       "walletsolidity/getexchangebyid?id=1",
		       "walletsolidity/listexchanges",
	        }
var eventQueryUrl =
	[]string{
		"blocks/total",
		"blocks",
		"blocks/latestSolidifiedBlockNumber",
		"transactions/total",
		"transactions",
		"transfers/total",
		"transfers",
		"events",
		"events/total",
		"events/timestamp",
		"events/TPt8DTDBZYfJ9fuyRjdWJr4PP68tRfptLG",
		"events/transaction/381df46287c296937582863269836dba8b6fc2098247fd86c2467ec7395ea854",
		"trc20/getholder/TPt8DTDBZYfJ9fuyRjdWJr4PP68tRfptLG",
		"contractlogs",
		"contractlogs/transaction/381df46287c296937582863269836dba8b6fc2098247fd86c2467ec7395ea854",
		"contractlogs/contract/TPt8DTDBZYfJ9fuyRjdWJr4PP68tRfptLG",
		"contractlogs/total",
	}
var eventQueryServiceIp = []string {
	"https://api.tronex.io",
	"http://47.254.69.13:8080",
	"http://47.254.84.58:8080",
}

var tronStackIp = []string {
	"https://api.tronstack.io",
	"http://18.221.34.0:8090",
	"http://52.15.93.92:8090",
}

func main() {
	logs.Info("start monitor")
	go start()
	go report()
	go change()
	go httpReport()
	go maxBlockReportAlert()
	defer influxdb.Client.C.Close()

	beego.Run()
}

func maxBlockReportAlert()  {
	c := cron.New()
	c.AddFunc("0 0,10,20,30,40,50 * * * *", func() {
		alerts.MaxBlockReportAlert()
	})
	c.Start()
}

func monitorEventQuery(getNowBlockAlert *alerts.GetNowBlockAlert) {
	tronStackAlert := new(tronstack.TronStackAlert)
	getNowBlockAlert.ReportEventQuery(httpevent, eventQueryUrl)
	for _, ip := range eventQueryServiceIp {
		tronStackAlert.CheckEventQueryStatus(ip)
	}

	for _, ip := range tronStackIp {
		tronStackAlert.CheckTronStackStatus(ip)
	}
	tronStackAlert.CheckTronStackResponse(eventQueryServiceIp[0], eventQueryServiceIp[1:])
}

func httpReport() {
	c := cron.New()
	c.AddFunc("0,20,40 * * * * *", func() {
		getNowBlockAlert := new(alerts.GetNowBlockAlert)
		getNowBlockAlert.Load()
		//getNowBlockAlert.ReportSRDelay()
		getNowBlockAlert.ReportDelayEX(httpex, fullUrl, solidityUrl)
		monitorEventQuery(getNowBlockAlert)
		getNowBlockAlert.ReportDelay(fullUrl, solidityUrl)
		//getNowBlockAlert.ReportEventQuery(httpTestEvent, eventQueryUrl)
	})
	c.Start()
}

func callWhoOnDuty() string {
	DUTY := []string{"吴彬", "岳瑞鹏", "张思聪", "吴斌", "梁志彦", "张博", "成特学"}
	who  := (time.Now().Unix() - 86400 * 4) / 86400 / 7 % int64(len(DUTY))
	return DUTY[who]
}

func report() {
	c := cron.New()
	c.AddFunc("0 0 2 * * *", func() {
		logs.Debug("report start")
		r := new(reports.TotalMissed)
		r.Date = time.Now().AddDate(0, 0, -1)
		logs.Debug("report date", r.Date.Format("2006-01-02 15:04:05"))
		r.ComputeData()
		r.Save()
		r.Report()
		context := "今天由" + callWhoOnDuty()+ "负责主网所有问题"
		slack.SendSlackNotification(context)
	})
	c.Start()
}

func change() {
	c := new(alerts.ChainParameters)
	c.MonitorUrl = config.MonitorConfig.Task.ProposalsMonitorUrl
	logs.Info("init proposals monitor url:", c.MonitorUrl)

	ticker := time.NewTicker(
		time.Duration(config.MonitorConfig.Task.GetDataInterval) *
			time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.RequestData()
			c.Judge()
		}
	}
}

func start() {
	for _, r := range datamanger.Requests {
		r.Load()
	}

	getNowBlockAlert := new(alerts.GetNowBlockAlert)
	getNowBlockAlert.Load()

	listWitnessAlert := new(alerts.ListWitnessesAlert)
	listWitnessAlert.Load()

	ticker := time.NewTicker(
		time.Duration(config.MonitorConfig.Task.GetDataInterval) *
			time.Second)
	defer ticker.Stop()

	startAlertCount := 0
	alertFinish := true

	for {
		select {
		case <-ticker.C:
			logs.Debug("start")

			for _, r := range datamanger.Requests {
				go r.Request()
			}

			time.Sleep(10 * time.Second)
			startAlertCount++

			if startAlertCount > 10 && alertFinish {
				alertFinish = false
				getNowBlockAlert.Start()
				getNowBlockAlert.Alert()

				listWitnessAlert.Start()
				listWitnessAlert.Alert()
				alertFinish = true
			}
		}
	}
}

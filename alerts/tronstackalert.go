package alerts

import (
	"github.com/sasaxie/monitor/common/database/influxdb"
	"github.com/sasaxie/monitor/function"
	"time"
)


type TronStackAlert struct {
}

func (alert TronStackAlert) CheckHeathStatus(ip string) {
	healthCheck := ip + "/healthcheck"
	status := function.CheckHeathcheck(healthCheck)
	httpMap := map[string]string{
		"IP":   ip,
	}
	httpFields := map[string]interface{}{
		"Status":   status,
		"IP":      ip,
	}
	influxdb.Client.WriteByTime("api_report_eventQuery_status", httpMap, httpFields, time.Now())
}


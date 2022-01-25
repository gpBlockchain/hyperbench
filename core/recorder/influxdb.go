package recorder

import (
	"github.com/influxdata/influxdb/client"
	"github.com/meshplus/hyperbench/common"
	"net/url"
	"time"
)

type BlockChainConfig struct {
	blockChainType string
	benchmarkName string
	config string
}

type influxdb struct {
	url       *url.URL
	database  string
	username  string
	password  string
	blkConfig BlockChainConfig

	client *client.Client
}

func (i *influxdb) process(report common.Report) {
	go i.sendProcess(report)
}

func (i *influxdb) sendProcess(report common.Report) {
	pts := make([]client.Point, 0, len(report.Cur.Results))

	for _, r := range report.Cur.Results {
		pts = append(pts, client.Point{
			Measurement: "current",
			Tags: map[string]string{
				"label": r.Label,
			},
			Time: time.Unix(0, r.Time),
			Fields: map[string]interface{}{
				"send":         r.Num,
				"duration":     r.Duration,
				"succ":         r.Statuses[common.Success],
				"fail":         r.Statuses[common.Failure],
				"confirmed":    r.Statuses[common.Confirm],
				"unknown":      r.Statuses[common.Unknown],
				"send_avg":     r.Send.Avg,
				"send_p0":      r.Send.P0,
				"send_p50":     r.Send.P50,
				"send_p90":     r.Send.P90,
				"send_p95":     r.Send.P95,
				"send_p99":     r.Send.P99,
				"send_p100":    r.Send.P100,
				"confirm_avg":  r.Confirm.Avg,
				"confirm_p0":   r.Confirm.P0,
				"confirm_p50":  r.Confirm.P50,
				"confirm_p90":  r.Confirm.P90,
				"confirm_p95":  r.Confirm.P95,
				"confirm_p99":  r.Confirm.P99,
				"confirm_p100": r.Confirm.P100,
				"write_avg":    r.Write.Avg,
				"write_p0":     r.Write.P0,
				"write_p50":    r.Write.P50,
				"write_p90":    r.Write.P90,
				"write_p95":    r.Write.P95,
				"write_p99":    r.Write.P99,
				"write_p100":   r.Write.P100,
			},
		})
	}

	bps := client.BatchPoints{
		Points:   pts,
		Database: i.database,
	}

	_, err := i.client.Write(bps)
	if err != nil {
		common.GetLogger("infx").Error(err)
	}
}



func (i *influxdb)statistic(statistic common.RemoteStatistic){
	go i.sendStatistic(statistic)
}



func (i *influxdb) sendStatistic(statistic common.RemoteStatistic) {
	// base info
	pts := make([]client.Point, 0, 1)

	pts = append(pts, client.Point{
		Measurement: "statistic",
		Tags: map[string]string{
			"label": "label",
		},
		Time: time.Unix(0, statistic.Start),
		Fields: map[string]interface{}{
			"blockChainType":i.blkConfig.blockChainType,
			"benchmarkName":i.blkConfig.benchmarkName,
			"config":i.blkConfig.config,
			"Start":statistic.Start,
			"End":statistic.End,
			"Tps":statistic.Tpss,
			"AvgTps":statistic.Tpss,
			"Bps":statistic.Bpss,
			"TxNum":statistic.TxNum,
			"BlockNum":statistic.BlockNum,
		},
	})

	bps := client.BatchPoints{
		Points:   pts,
		Database: i.database,
	}

	_, err := i.client.Write(bps)
	if err != nil {
		common.GetLogger("infx").Error(err)
	}
}


func (i *influxdb) release(){
}



func newInfluxdb(blkConfig BlockChainConfig, URL string, database string, username string, password string) (*influxdb, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	i := &influxdb{
		url:       u,
		database:  database,
		username:  username,
		password:  password,
		blkConfig: blkConfig,

	}
	err = i.makeClient()
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *influxdb) makeClient() (err error) {
	i.client, err = client.NewClient(client.Config{
		URL:      *i.url,
		Username: i.username,
		Password: i.password,
		Timeout:  30 * time.Second,
	})
	return
}

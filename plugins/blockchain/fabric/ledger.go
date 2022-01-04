package fabric

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/meshplus/hyperbench/common"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"
	"time"
)

// GetTPS get remote tps
func GetTPS(client *ledger.Client, startNum uint64, statistic bcom.Statistic) (*common.RemoteStatistic, error) {

	blockInfo, err := client.QueryInfo()
	if err != nil {
		return nil, err
	}
	height := blockInfo.BCI.Height
	var (
		cur           uint64
		txNum         int
		totalTxNum    int
		blockNum      int
		totalBlockNum int
		timeSeq       []int64
		heightSeq     []uint64
		maxTps        int
		maxBps        int
		tpss          []int
		bpss          []int
		from          int64
	)
	// add last query info, nanoseconds to seconds
	timeSeq = statistic.TimeSeq
	timeSeq = append(timeSeq, statistic.To)
	heightSeq = statistic.HeightSeq
	heightSeq = append(heightSeq, height)
	from = statistic.From
	maxBps = 0
	maxBps = 0
	for i := 0; i < len(heightSeq); i++ {
		txNum = 0
		blockNum = 0
		if cur > heightSeq[i] {
			// skip error query result
			continue
		}
		for cur = startNum; cur < heightSeq[i]; cur++ {
			block, err := client.QueryBlock(cur)
			if err != nil {
				return nil, err
			}
			txNum += len(block.GetData().Data)
			blockNum++
		}
		// check for divide by 0
		interval := int(time.Unix(0, timeSeq[i]).Sub(time.Unix(0, from)).Seconds())
		if interval <= 0 {
			interval = 1
		}
		tps := txNum / interval
		if tps > maxTps {
			maxTps = tps
		}
		bps := blockNum / interval
		if bps > maxBps {
			maxBps = bps
		}
		tpss = append(tpss, tps)
		bpss = append(bpss, bps)

		totalTxNum += txNum
		totalBlockNum += blockNum
		startNum = heightSeq[i]
		from = timeSeq[i]
	}
	tpss = append(tpss, maxTps)
	bpss = append(bpss, maxBps)

	statisticResult := &common.RemoteStatistic{
		Start:    statistic.From,
		End:      statistic.To,
		BlockNum: totalBlockNum,
		TxNum:    totalTxNum,
		Tpss:     tpss,
		Bpss:     bpss,
	}
	return statisticResult, nil
}

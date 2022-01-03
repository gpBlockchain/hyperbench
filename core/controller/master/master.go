package master

import (
	"time"

	"github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/vm"
	"github.com/meshplus/hyperbench/vm/base"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

// Master is the interface of mater node
type Master interface {
	// Prepare is used to prepare
	Prepare() error

	// GetContext generate the context, which will be passed to Worker
	GetContext() ([]byte, error)

	// Statistic query the remote statistic data from chain
	Statistic(from, to int64) (*common.RemoteStatistic, error)

	// log height and time
	LogBlockHeightWithTime()

	// GetTpsUnit return the time interval for caculating tps
	GetTpsUnit() time.Duration
}

// LocalMaster is the implement of master in local
type LocalMaster struct {
	masterVM       vm.VM
	heightLogs     *heightWithTimeSeq
	params         []string
	tpsUnit        time.Duration
	engineDuration time.Duration
}

type heightWithTimeSeq struct {
	timeSeq   []int64
	heightSeq []uint64
}

// Prepare is used to prepare
func (m *LocalMaster) Prepare() (err error) {
	// call user hook
	err = m.masterVM.BeforeDeploy()
	if err != nil {
		return errors.Wrap(err, "can not call user hook `BeforeDeploy`")
	}

	// prepare contract
	err = m.masterVM.DeployContract()
	if err != nil {
		return errors.Wrap(err, "can not deploy contract")
	}

	return nil
}

// GetContext generate the context, which will be passed to Worker
func (m *LocalMaster) GetContext() ([]byte, error) {
	err := m.masterVM.BeforeGet()
	if err != nil {
		return nil, err
	}
	return m.masterVM.GetContext()
}

// Statistic query the remote statistic data from chain
func (m *LocalMaster) Statistic(from, to int64) (*common.RemoteStatistic, error) {
	return m.masterVM.Statistic(from, to, m.heightLogs.timeSeq, m.heightLogs.heightSeq)
}

// LogBlockHeightWithTime record timestamp and chain height
func (m *LocalMaster) LogBlockHeightWithTime() {
	queryTime, height := m.masterVM.LogLedgerHeight()
	m.heightLogs.timeSeq = append(m.heightLogs.timeSeq, queryTime)
	m.heightLogs.heightSeq = append(m.heightLogs.heightSeq, height)
}

// GetTpsUnit return the time interval for caculating tps
func (m *LocalMaster) GetTpsUnit() time.Duration {
	if m.tpsUnit.Seconds() > m.engineDuration.Seconds() {
		m.tpsUnit = m.engineDuration
	}
	return m.tpsUnit
}

// NewLocalMaster create LocalMaster.
func NewLocalMaster() (*LocalMaster, error) {

	params := viper.GetStringSlice(common.ClientContractArgsPath)
	scriptPath := viper.GetString(common.ClientScriptPath)
	tpsUnit := viper.GetDuration(common.EngineTpsUnitPath)
	engineDuration := viper.GetDuration(common.EngineDurationPath)
	vmType := strings.TrimPrefix(filepath.Ext(scriptPath), ".")
	masterVM, err := vm.NewVM(vmType, base.ConfigBase{
		Path: scriptPath,
		Ctx: common.VMContext{
			WorkerIdx: -1,
			VMIdx:     -1,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "can not create master")
	}
	ht := &heightWithTimeSeq{
		timeSeq:   make([]int64, 0),
		heightSeq: make([]uint64, 0),
	}
	return &LocalMaster{
		masterVM:       masterVM,
		params:         params,
		heightLogs:     ht,
		tpsUnit:        tpsUnit,
		engineDuration: engineDuration,
	}, nil
}

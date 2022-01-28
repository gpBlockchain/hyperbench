package glua

import (
	"fmt"
	"github.com/meshplus/hyperbench-common/base"
	fcom "github.com/meshplus/hyperbench-common/common"
)

type FakeChain struct {
	Name string
	base *base.BlockchainBase
}

func NewMock() (client *FakeChain, err error) {
	return &FakeChain{"fake",
		base.NewBlockchainBase(base.ClientConfig{}),
	}, nil
}

func (chain *FakeChain) DeployContract() error {
	return nil
}

func (chain *FakeChain) Invoke(invoke fcom.Invoke, ops ...fcom.Option) *fcom.Result {
	fmt.Printf("invoke:%v\n", invoke)
	fmt.Printf("ops:%v\n", ops)
	return &fcom.Result{Status: fcom.Success}
}

func (chain *FakeChain) Transfer(fcom.Transfer, ...fcom.Option) *fcom.Result {
	return &fcom.Result{}
}

func (chain *FakeChain) Confirm(rt *fcom.Result, ops ...fcom.Option) *fcom.Result {
	return &fcom.Result{}
}

func (chain *FakeChain) Query(bq fcom.Query, ops ...fcom.Option) interface{} {
	return nil
}

func (chain *FakeChain) Option(fcom.Option) error {
	return nil
}

func (chain *FakeChain) GetContext() (string, error) {
	return "", nil
}

func (chain *FakeChain) SetContext(ctx string) error {
	return nil
}

func (chain *FakeChain) ResetContext() error {
	return nil
}

func (chain *FakeChain) Statistic(statistic fcom.Statistic) (*fcom.RemoteStatistic, error) {
	return &fcom.RemoteStatistic{}, nil
}

func (chain *FakeChain) LogStatus() (int64, error) {
	return 0, nil
}

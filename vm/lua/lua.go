// lua implement of vm.VM
package lua

import (
	"errors"

	"github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/plugins/blockchain"
	base2 "github.com/meshplus/hyperbench/plugins/blockchain/base"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"
	idex "github.com/meshplus/hyperbench/plugins/index"
	"github.com/meshplus/hyperbench/plugins/toolkit"
	"github.com/meshplus/hyperbench/vm/base"
	"github.com/spf13/viper"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// VM the implementation of BaseVM for lua.
type VM struct {
	*base.VMBase

	vm       *lua.LState
	instance *lua.LTable
	meta     *lua.LTable
	client   blockchain.Blockchain

	index *idex.Index
}

// NewVM use given base to create VM.
func NewVM(base *base.VMBase) (vm *VM, err error) {

	vm = &VM{
		VMBase: base,
		index: &idex.Index{
			Worker: base.Ctx.WorkerIdx,
			VM:     base.Ctx.VMIdx,
		},
	}

	vm.vm = lua.NewState()
	defer func() {
		if err != nil {
			vm = nil
		}
	}()

	// inject test case metatable
	vm.injectTestcaseBase()

	// append plugins to the test case metatable
	err = vm.setPlugins(vm.meta)
	if err != nil {
		return nil, err
	}

	// load script
	err = vm.vm.DoFile(base.Path)
	if err != nil {
		return nil, err
	}

	// get test
	var ok bool
	vm.instance, ok = vm.vm.Get(-1).(*lua.LTable)
	if !ok {
		return nil, errors.New("script's return value is not table")
	}
	vm.vm.Pop(1)

	return vm, nil
}

// hooks
const (
	beforeDeploy = "BeforeDeploy"
	// nolint
	deployContract = "DeployContract"
	beforeGet      = "BeforeGet"
	beforeSet      = "BeforeSet"
	// nolint
	setContext = "SetContext"
	beforeRun  = "BeforeRun"
	run        = "Run"
	afterRun   = "AfterRun"
)

// builtin
const (
	lNew   = "new"
	lIndex = "__index"
)

// plugins
const (
	testcase = "testcase"
	client   = "blockchain"
	tool     = "toolkit"
	index    = "index"
)

func (v *VM) injectTestcaseBase() {
	mt := v.vm.NewTypeMetatable(testcase)
	v.vm.SetGlobal(testcase, mt)

	var empty lua.LGFunction = func(state *lua.LState) int {
		return 0
	}
	var result lua.LGFunction = func(state *lua.LState) int {
		state.Push(luar.New(state, &common.Result{}))
		return 1
	}

	v.vm.SetField(mt, lNew, v.vm.NewFunction(func(state *lua.LState) int {
		table := v.vm.NewTable()
		v.vm.SetMetatable(table, v.vm.GetMetatable(lua.LString(testcase)))
		_ = v.setPlugins(table)
		v.vm.Push(table)
		return 1
	}))

	v.vm.SetField(mt, lIndex, v.vm.SetFuncs(v.vm.NewTable(), map[string]lua.LGFunction{
		beforeDeploy: empty,
		beforeGet:    empty,
		beforeRun:    empty,
		run:          result,
		afterRun:     empty,
	}))
	v.meta = mt
}

// BeforeDeploy will call before deploy contract.
func (v *VM) BeforeDeploy() error {
	fn := v.instance.RawGetString(beforeDeploy)
	if fn != lua.LNil {
		return v.vm.CallByParam(lua.P{
			Fn: fn,
		}, v.instance)
	}
	return nil
}

// DeployContract deploy contract.
func (v *VM) DeployContract() error {
	return v.client.DeployContract()
}

// BeforeGet will call before get context.
func (v *VM) BeforeGet() error {
	fn := v.instance.RawGetString(beforeGet)
	if fn != lua.LNil {
		return v.vm.CallByParam(lua.P{
			Fn: fn,
		}, v.instance)
	}
	return nil
}

// GetContext generate context for execute tx in vm.
func (v *VM) GetContext() ([]byte, error) {
	s, err := v.client.GetContext()
	return []byte(s), err
}

// Statistic statistic remote execute info.
func (v *VM) Statistic(from, to int64) (*common.RemoteStatistic, error) {

	return v.client.Statistic(bcom.Statistic{
		From: from,
		To:   to,
	})
}

// BeforeSet will call before set context.
func (v *VM) BeforeSet() error {
	fn := v.instance.RawGetString(beforeSet)
	if fn != lua.LNil {
		return v.vm.CallByParam(lua.P{
			Fn: fn,
		}, v.instance)
	}
	return nil
}

// SetContext set context for execute tx in vm, the ctx is generated by GetContext.
func (v *VM) SetContext(ctx []byte) error {
	return v.client.SetContext(string(ctx))
}

// BeforeRun will call once before run.
func (v *VM) BeforeRun() error {
	fn := v.instance.RawGetString(beforeRun)
	if fn != lua.LNil {
		return v.vm.CallByParam(lua.P{
			Fn: fn,
		}, v.instance)
	}
	return nil
}

// Run create and send tx to client.
func (v *VM) Run(ctx common.TxContext) (*common.Result, error) {
	v.index.Engine = ctx.EngineIdx
	v.index.Tx = ctx.TxIdx

	err := v.vm.CallByParam(lua.P{
		Fn:      v.instance.RawGetString(run),
		NRet:    1,
		Protect: false,
	}, v.instance)

	if err != nil {
		v.Logger.Error(err)
		return nil, err
	}
	val := v.vm.Get(-1)
	v.vm.Pop(1)

	ud, ok := val.(*lua.LUserData)
	if !ok {
		return nil, errors.New("returned val is not user data")
	}
	res, ok := ud.Value.(*common.Result)
	if !ok {
		return nil, errors.New("returned user data is not result")
	}
	return res, nil
}

// AfterRun will call once after run.
func (v *VM) AfterRun() error {
	fn := v.instance.RawGetString(afterRun)
	if fn != lua.LNil {
		return v.vm.CallByParam(lua.P{
			Fn: fn,
		}, v.instance)
	}
	return nil
}

// Close close vm.
func (v *VM) Close() {
	v.vm.Close()
}

func (v *VM) setPlugins(table *lua.LTable) (err error) {

	clientType, clientConfigPath := viper.GetString(common.ClientTypePath), viper.GetString(common.ClientConfigPath)
	options := viper.GetStringMap(common.ClientOptionPath)
	contractPath := viper.GetString(common.ClientContractPath)
	args, _ := viper.Get(common.ClientContractArgsPath).([]interface{})
	if clientType == "eth" {
		options["vmIdx"] = v.index.VM
		options["wkIdx"] = v.index.Worker
	}
	v.client, err = blockchain.NewBlockchain(base2.ClientConfig{
		ClientType:   clientType,
		ConfigPath:   clientConfigPath,
		ContractPath: contractPath,
		Args:         args,
		Options:      options,
	})

	if err != nil {
		return err
	}

	// todo: register the plugins manually instead of luar's reflection to optimize performance
	lClient := luar.New(v.vm, v.client)
	lToolKit := luar.New(v.vm, toolkit.NewToolKit())
	lIndex := luar.New(v.vm, v.index)
	v.vm.SetField(table, client, lClient)
	v.vm.SetField(table, tool, lToolKit)
	v.vm.SetField(table, index, lIndex)

	return nil
}

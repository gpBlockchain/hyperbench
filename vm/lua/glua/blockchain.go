package glua

import (
	"github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/plugins/blockchain"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"
	lua "github.com/yuin/gopher-lua"
)

func newBlockchain(L *lua.LState, client blockchain.Blockchain) lua.LValue {
	clientTable := L.NewTable()
	clientTable.RawSetString("DeployContract", DeployContractLuaFunction(L, client))
	clientTable.RawSetString("Invoke", InvokeLuaFunction(L, client))
	clientTable.RawSetString("Transfer", TransferLuaFunction(L, client))
	clientTable.RawSetString("Confirm", ConfirmLuaFunction(L, client))
	clientTable.RawSetString("Query", QueryLuaFunction(L, client))
	clientTable.RawSetString("Option", OptionLuaFunction(L, client))
	//todo support context
	//clientTable.RawSetString("GetContext",nil)
	//clientTable.RawSetString("SetContext",nil)
	//clientTable.RawSetString("ResetContext",nil)
	//clientTable.RawSetString("Statistic",nil)
	return clientTable
}

func OptionLuaFunction(L *lua.LState, client blockchain.Blockchain) lua.LValue {
	return L.NewFunction(func(state *lua.LState) int {
		var map1 bcom.Option
		// case.blockchain:Invoke() --> first arg is blockchain.Blockchain
		// case.blockchain.Invoke  ----> first arg is normal
		firstArgIndex := 1
		// check first arg is blockchain.Blockchain
		if checkBlockChainByIdx(state, 1) {
			firstArgIndex++
		}
		invokeTable := state.CheckTable(firstArgIndex)
		err := TableLua2GoStruct(invokeTable, &map1)
		if err != nil {
			state.ArgError(1, "common.Option expected")
		}
		err = client.Option(map1)
		if err != nil {
			state.Push(lua.LString(err.Error()))
		}
		state.Push(lua.LString(""))
		return 1
	})
}

func InvokeLuaFunction(L *lua.LState, client blockchain.Blockchain) *lua.LFunction {
	var invoke bcom.Invoke
	return invokeLuaFunction(L, client, invoke, func(b blockchain.Blockchain, b2 interface{}, option ...bcom.Option) interface{} {
		return b.Invoke(b2.(bcom.Invoke), option...)
	})
}

func TransferLuaFunction(L *lua.LState, client blockchain.Blockchain) *lua.LFunction {
	var transfer bcom.Transfer
	return invokeLuaFunction(L, client, transfer, func(b blockchain.Blockchain, b2 interface{}, option ...bcom.Option) interface{} {
		return b.Transfer(b2.(bcom.Transfer), option...)
	})
}

func QueryLuaFunction(L *lua.LState, client blockchain.Blockchain) *lua.LFunction {
	var query bcom.Query
	return invokeLuaFunction(L, client, query, func(b blockchain.Blockchain, b2 interface{}, option ...bcom.Option) interface{} {
		return b.Query(b2.(bcom.Query), option...)
	})
}

func ConfirmLuaFunction(L *lua.LState, client blockchain.Blockchain) *lua.LFunction {
	var confirm common.Result
	return invokeLuaFunction(L, client, confirm, func(b blockchain.Blockchain, b2 interface{}, option ...bcom.Option) interface{} {
		return b.Confirm(b2.(*common.Result), option...)
	})
}

func invokeLuaFunction(L *lua.LState, cli blockchain.Blockchain, arg1Type interface{}, fn func(blockchain.Blockchain, interface{}, ...bcom.Option) interface{}) *lua.LFunction {
	return L.NewFunction(func(state *lua.LState) int {
		// case.blockchain:Invoke() --> first arg is blockchain.Blockchain
		// case.blockchain.Invoke  ----> first arg is normal
		firstArgIndex := 1
		// check first arg is blockchain.Blockchain
		if checkBlockChainByIdx(state, 1) {
			firstArgIndex++
		}
		invokeTable := state.CheckTable(firstArgIndex)
		err := TableLua2GoStruct(invokeTable, &arg1Type)
		if err != nil {
			state.ArgError(1, "interface. expected")
		}
		if state.GetTop() == 1+firstArgIndex {
			ret := fn(cli, arg1Type)
			state.Push(go2Lua(state, ret))
			return 1
		}
		var opts []bcom.Option
		for i := 1 + firstArgIndex; i <= state.GetTop(); i++ {
			table := state.CheckTable(i)
			var map1 bcom.Option
			err := TableLua2GoStruct(table, &map1)
			if err != nil {
				state.ArgError(1, "common.Option expected")
			}
			opts = append(opts, map1)
		}
		ret := fn(cli, arg1Type, opts...)
		state.Push(go2Lua(state, ret))
		return 1
	})
}

func checkBlockChainByIdx(state *lua.LState, idx int) bool {
	if state.GetTop() > idx {
		return false
	}
	lvalue := state.CheckTable(idx)
	k, _ := lvalue.Next(lua.LString("Invok"))
	// check arg is blockchain.Blockchain
	if k.String() != "Invoke" {
		return false
	}
	return true
}

func DeployContractLuaFunction(L *lua.LState, client blockchain.Blockchain) *lua.LFunction {
	return L.NewFunction(func(state *lua.LState) int {
		err := client.DeployContract()
		if err != nil {
			state.Push(lua.LString(err.Error()))
			return 1
		}
		state.Push(lua.LString(""))
		return 1
	})
}

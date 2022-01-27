package glua

import (
	"github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/plugins/blockchain"
	idex "github.com/meshplus/hyperbench/plugins/index"
	"github.com/meshplus/hyperbench/plugins/toolkit"
	lua "github.com/yuin/gopher-lua"
)

func NewClientLValue(L *lua.LState, client blockchain.Blockchain) lua.LValue {
	return newBlockchainTable(L, client)
}

func NewToolKitLValue(L *lua.LState, kit *toolkit.ToolKit) lua.LValue {
	return newToolKitTable(L, kit)
}

func NewLIndexLValue(L *lua.LState, idx *idex.Index) lua.LValue {
	return newIdexIndexTable(L, idx)
}

func NewResultLValue(L *lua.LState, r *common.Result) lua.LValue {
	return newCommonResultTable(L, r)
}

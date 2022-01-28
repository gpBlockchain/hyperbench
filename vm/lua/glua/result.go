package glua

import (
	"github.com/meshplus/hyperbench/common"
	lua "github.com/yuin/gopher-lua"
)

func newCommonResult(L *lua.LState, r *common.Result) lua.LValue {
	//todo replace reflect
	resultTable := L.NewTable()
	resultTable.RawSetString("Label", lua.LString(r.Label))
	resultTable.RawSetString("UID", lua.LString(r.UID))
	resultTable.RawSetString("BuildTime", lua.LNumber(r.BuildTime))
	resultTable.RawSetString("SendTime", lua.LNumber(r.SendTime))
	resultTable.RawSetString("WriteTime", lua.LNumber(r.WriteTime))
	resultTable.RawSetString("Status", lua.LString(r.Status))
	//todo support ret
	//resultTable.RawSetString("Status",lua.LTable{}(r.Status))
	return resultTable
}

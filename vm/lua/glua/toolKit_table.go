package glua

import (
	"github.com/meshplus/hyperbench/plugins/toolkit"
	lua "github.com/yuin/gopher-lua"
)



func HexLuaFunction(L *lua.LState,kit *toolkit.ToolKit) *lua.LFunction{
	return L.NewFunction(func(state *lua.LState) int {
		input := state.CheckString(1)
		ret := kit.Hex(input)
		state.Push(lua.LString(ret))
		return 1
	})
}

func randStrLuaFunction(L *lua.LState,kit *toolkit.ToolKit) *lua.LFunction {
	return L.NewFunction(func(state *lua.LState) int {
		size := state.CheckInt(1)
		ret := kit.RandStr(uint(size))
		L.Push(lua.LString(ret))
		return 1
	})
}

func StringLuaFunction(L *lua.LState,kit *toolkit.ToolKit) *lua.LFunction{
	return L.NewFunction(func(state *lua.LState) int {
		argLength := state.GetTop()
		if argLength <2{
			panic("args are less than 1")
		}
		input := state.CheckAny(1)
		if argLength == 1{
			ret := kit.String(input)
			L.Push(lua.LString(ret))
			return 1
		}
		var offsets []int
		for i := 2; i < argLength; i++ {
			offset := state.CheckInt(i)
			offsets = append(offsets, offset)
		}
		ret := kit.String(input,offsets...)
		L.Push(lua.LString(ret))
		return 1
	})
}

func TestInterfaceLuaFunction(L *lua.LState,kit *toolkit.ToolKit) *lua.LFunction{
	return L.NewFunction(func(state *lua.LState) int {
		input := state.CheckAny(1)
		arg :=ToGoValue(input,Option{})
		ret := kit.TestInterface(arg)
		state.Push(Go2Lua(state, ret))
		return 1
	})
}


func newToolKitTable(L *lua.LState, kit *toolkit.ToolKit) lua.LValue {
	toolkitTable := L.NewTable()
	toolkitTable.RawSetString("Hex",HexLuaFunction(L,kit))
	toolkitTable.RawSetString("RandStr",randStrLuaFunction(L,kit))
	toolkitTable.RawSetString("String",StringLuaFunction(L,kit))
	toolkitTable.RawSetString("TestInterface",TestInterfaceLuaFunction(L,kit))
	return toolkitTable
}







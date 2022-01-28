package glua

import (
	"fmt"
	"github.com/meshplus/hyperbench/plugins/toolkit"
	lua "github.com/yuin/gopher-lua"
	"reflect"
	"testing"
)

var script1 = `
	print("-----1-----")
	print(u.toolkit:RandStr(10))
	print(u.toolkit:randStr(10))
	print("----2-----")
	print(u.toolkit.Name)
	print("----3-----")
	u.toolkit.Name = "hhe"
	print(u.toolkit.Name)
	print(u.toolkit:Hex("aaaaaa"))
	-- print(u.toolkit:String((1,2,3,4),1,2))
	print("TestInterfaceLua",u.toolkit:TestInterface({aa=bb}))
`

var script2 = `
	print("-----1-----")
	print(u.toolkit.RandStr(10))
	print(u.toolkit.Hex("aaaaaa"))
	-- print(u.toolkit.String((1,2,3,4),1,2))
	print("TestInterfaceLua",u.toolkit.TestInterface({Aa=Bb}))
`

func Test_toolkit(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("u")
	L.SetGlobal("u", mt)
	L.SetField(mt, "toolkit", newToolKit(L, toolkit.NewToolKit()))
	if err := L.DoString(script1); err != nil {
		panic(err)
	}
}

func Test_toolkit_table(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("u")
	L.SetGlobal("u", mt)
	L.SetField(mt, "toolkit", newToolKit(L, toolkit.NewToolKit()))
	if err := L.DoString(script2); err != nil {
		panic(err)
	}
}

//55 :45
func BenchmarkName(b *testing.B) {
	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("u")
	L.SetGlobal("u", mt)
	L.SetField(mt, "toolkit", newToolKit(L, &toolkit.ToolKit{}))
	L.SetField(mt, "toolkitMetatable", newTookKit(L))
	scriptTable := `
	u.toolkit.Hex("aaaaaa")
	`
	scriptLua := `
	u.toolkitMetatable:Hex("aaaaaa")
	`
	for i := 0; i < b.N; i++ {
		runMetaTable(L, scriptLua)
		runTable(L, scriptTable)
	}
}
func runMetaTable(L *lua.LState, script string) {
	if err := L.DoString(script); err != nil {
		panic(err)
	}
}
func runTable(L *lua.LState, script string) {
	if err := L.DoString(script); err != nil {
		panic(err)
	}
}

func Test_toolkit2(t *testing.T) {
	kit := toolkit.NewToolKit()
	fmt.Println("String:", kit.String([]byte{1, 2, 3, 4}, 0, 3))

}

func Test_1234(t *testing.T) {
	var value interface{}
	value = "1"
	//ret := value.(type)
	//fmt.Println(ret)
	switch converted := value.(type) {
	case bool:
		fmt.Println(converted)
	case int:
		fmt.Println(converted)
	case string:
		fmt.Println("str:", converted)

	}
}

func Test_12345(t *testing.T) {
	L := lua.NewState()
	t1 := L.NewTable()
	t2 := L.NewTable()
	t1.RawSetString("key", lua.LNumber(111))
	t1.RawSetInt(2, lua.LNumber(222))
	t1.RawSetString("func", L.NewFunction(func(state *lua.LState) int {
		L.Push(lua.LString("heheeh"))
		//L.Get()
		return 1
	}))

	t1.RawSetInt(4, t2) // 模拟map
	mt := L.NewTypeMetatable("u")
	L.SetGlobal("u", mt)
	L.SetField(mt, "demo", t1)
	var script1 = `
		print("-----1-----")
		print(u.demo.key)
		u.demo.key=2
		print(u.demo.key)
		print(u.demo.func())
	`
	if err := L.DoString(script1); err != nil {
		panic(err)
	}
}

func Test_123452(t *testing.T) {
	kit := toolkit.ToolKit{}
	fmt.Println("get kit 成员变量")
	getFieldByReflect(kit)
	fmt.Println("get kit 方法名 和参数")
	fmt.Println("通过方法名调用kit的方法")

}

func getFieldByReflect(data interface{}) interface{} {
	typeOfData := reflect.TypeOf(data)
	valueOfData := reflect.ValueOf(data)
	switch typeOfData.Kind() {
	case reflect.Struct:
		fmt.Println("struct")
		var strs []interface{}
		for i := 0; i < typeOfData.NumField(); i++ {
			strs = append(strs, getFieldByReflect(valueOfData.Index(i).Interface()))
		}
		return strs
	case reflect.Ptr:
		fmt.Println("Ptr")
		return getFieldByReflect(valueOfData.Elem().Interface())
	default:
		fmt.Println(typeOfData.Kind().String())
		return typeOfData.Kind().String()
	}
}

func newTookKit(L *lua.LState) lua.LValue {
	return NewToolKitLValue(L, &toolkit.ToolKit{})
}

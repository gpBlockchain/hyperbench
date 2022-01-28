package glua

import (
	"fmt"
	"github.com/meshplus/hyperbench/plugins/blockchain/fake"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

var s1 = `
	 ret=case.blockchain:Invoke({
		   func="123",
		   args={"123", "123"}
		},{aa="aa"},{bb="bb"})
`
var s2 = `
	 ret=case.blockchain.Invoke({
		   func="123",
		   args={"123", "123"}
		},{aa="aa"},{bb="bb"})
`

func Test_client(t *testing.T) {

	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("case")
	L.SetGlobal("case", mt)
	client, _ := fake.New()
	cLua := newBlockchain(L, client)
	L.SetField(mt, "blockchain", cLua)
	if err := L.DoString(s1); err != nil {
		panic(err)
	}
	fmt.Println("----s2--------")
	if err := L.DoString(s2); err != nil {
		panic(err)
	}
}


func Test_client_table(t *testing.T) {

	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("case")
	L.SetGlobal("case", mt)
	client, _ := fake.New()
	cLua := newBlockchain(L, client)
	L.SetField(mt, "blockchain", cLua)
	if err := L.DoString(s2); err != nil {
		panic(err)
	}
}

func luar3(l2 *lua.LState) {
	if err := l2.DoString(s2); err != nil {
		panic(err)
	}

}

func luar1(L *lua.LState) {
	if err := L.DoString(s1); err != nil {
		panic(err)
	}
}

func luar2(L *lua.LState) {
	if err := L.DoString(s1); err != nil {
		panic(err)
	}
}

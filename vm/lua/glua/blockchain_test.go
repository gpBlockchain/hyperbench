package glua

import (
	"github.com/meshplus/hyperbench/plugins/blockchain/fake"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
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
}


func Test_client_table(t *testing.T) {

	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("case")
	L.SetGlobal("case", mt)
	client, _ := fake.New()
	cLua := newBlockchainTable(L, client)
	L.SetField(mt, "blockchain", cLua)
	if err := L.DoString(s2); err != nil {
		panic(err)
	}
}

//handle BenchmarkInvokeLua-8   	   44344	     28033 ns/op
//auto   BenchmarkInvokeLua-8   	   41647	     26326 ns/op
func BenchmarkInvokeLua(b *testing.B) {
	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("case")
	L.SetGlobal("case", mt)
	client, _ := fake.New()
	cLua := luar.New(L, client)
	L.SetField(mt, "blockchain", cLua)

	L1 := lua.NewState()
	defer L1.Close()
	mt1 := L1.NewTypeMetatable("case")
	L1.SetGlobal("case", mt1)
	client1, _ := fake.New()
	cLua1 := newBlockchain(L1, client1)
	L1.SetField(mt1, "blockchain", cLua1)

	L2 := lua.NewState()
	defer L2.Close()
	mt2 := L2.NewTypeMetatable("case")
	L2.SetGlobal("case", mt2)
	client2, _ := fake.New()
	cLua2 := newBlockchainTable(L2, client2)
	L2.SetField(mt2, "blockchain", cLua2)

	for i := 0; i < b.N; i++ {
		luar1(L)
		luar2(L1)
		luar3(L2)
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

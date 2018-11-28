package main

import (
	"fmt"
	"3rdparty/src/glua"
)


type Int struct {
	I int
}

func NewInt() *Int {
	return &Int{10}
}

func (i Int) PrintInt(str string) {
	fmt.Println(str, i.I)
}

func main() {
	L := glua.NewState()

	var tlib = glua.Libfuncs{
		"gotest", // lib name
		map[string]interface{}{
			"NewInt":    NewInt,          // lua function name, go function
			"PrintInt":  (*Int).PrintInt, // lua function name, go function
			"goprintln": fmt.Println,
		},
	}
	if ok, err := L.Register(&tlib); !ok {
		println(err.Error())
	}
	L.Dostring("gotest.PrintInt(gotest.NewInt(), \"Int is\")")
	L.Dostring("gotest.goprintln(true, 123, \"lua\", gotest.NewInt())")
}
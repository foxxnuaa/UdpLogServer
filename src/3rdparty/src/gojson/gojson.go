package gojson

/*
#cgo CFLAGS: -I../golua/lua/lua
#cgo windows LDFLAGS:-L../golua/lua -llua_windows
#cgo llua LDFLAGS: -llua
#cgo linux LDFLAGS:-L../golua/lua -lluajit_linux -lm -ldl
#cgo darwin LDFLAGS:  -llua_osx
#cgo freebsd LDFLAGS: -llua

#include <lua.h>
#include <stdlib.h>

int luaopen_cjson(lua_State *l);

*/
import "C"

import (
	lua "3rdparty/src/golua/lua"
	"unsafe"
)

func Reg(L *lua.State) {
	C.luaopen_cjson((*C.lua_State)(unsafe.Pointer(L.GetCState())))
	L.SetGlobal("cjson")
}

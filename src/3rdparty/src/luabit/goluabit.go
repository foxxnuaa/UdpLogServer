package luabit

/*
#cgo CFLAGS: -I../golua/lua/lua
#cgo windows LDFLAGS:-L../golua/lua -llua_windows
#cgo linux LDFLAGS:-L../golua/lua -lluajit_linux -lm -ldl
#cgo darwin LDFLAGS: -llua_osx
#cgo freebsd LDFLAGS:-L../golua/lua -llua

#include <lua.h>
#include <stdlib.h>

int luaopen_bit(lua_State *L);

*/
import "C"

import (
	lua "3rdparty/src/golua/lua"
	"unsafe"
)

func Reg(L *lua.State) {
	C.luaopen_bit((*C.lua_State)(unsafe.Pointer(L.GetCState())))
	L.SetGlobal("bit")
}

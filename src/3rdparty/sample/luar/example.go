package main

import "fmt"
import "3rdparty/src/luar"

type MyStruct struct {
	Name string
	Age  int
}

const test = `
for i = 1,5 do
    Print(MSG,i)
end
Print(ST)
print(ST.Name,ST.Age)
--// slices!
for i,v in pairs(S) do
   print(i,v)
end
`

func main() {
	L := luar.Init()
	defer L.Close()

	S := []string{"alfred", "alice", "bob", "frodo"}

	ST := &MyStruct{"Dolly", 46}

	luar.Register(L, "", luar.Map{
		"Print": fmt.Println,
		"MSG":   "hello", // can also register constants
		"ST":    ST,      // and other values
		"S":     S,
	})

	L.DoString(test)

}

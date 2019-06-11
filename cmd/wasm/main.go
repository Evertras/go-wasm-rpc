//+build js

package main

import (
	"syscall/js"

	"github.com/Evertras/go-wasm-rpc/lib/sample"
)

func main() {
	c := make(chan bool, 0)

	js.Global().Set("gowasm", map[string]interface{}{})
	base := js.Global().Get("gowasm")

	sample.RegisterWasmCallbacks(base)

	println("WASM initialized")

	<-c
}

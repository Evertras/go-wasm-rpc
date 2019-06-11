package sample

import (
	"github.com/Evertras/go-wasm-rpc/lib/service"
)

type wasmService struct{}

// Make sure that we fit the proto-defined interface or fail fast at compile time.
// This doesn't have to be a separate package, but it makes the pipeline cleaner.
var _ service.WasmServiceServer = &wasmService{}

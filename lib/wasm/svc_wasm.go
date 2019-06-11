package wasm

import (
	"context"

	"github.com/Evertras/go-wasm-rpc/lib/sample"
)

type wasmServer struct{}

// Make sure that we fit the proto-defined interface or fail fast at compile time.
// This doesn't have to be a separate package, but it makes the pipeline cleaner.
var _ sample.WasmServiceServer = &wasmServer{}

func (s *wasmServer) Echo(ctx context.Context, req *sample.EchoRequest) (*sample.EchoResponse, error) {
	res := &sample.EchoResponse{}

	res.Text = req.Text + " (with an echo!)"

	return res, nil
}

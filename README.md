# go-wasm-rpc

Use gRPC service definitions to communicate between WASM and JS.
The intent is to allow you to use an existing toolset (gRPC, protobuf)
to allow a simple, type-safe API between JS (with Typescript) and
a Go WASM module.

This repo consists both of the tool itself and a sample for reference.
Run `make` to build everything.  Run `./sample-server` after `make` to
run a server that will serve the front end on `localhost:8000`, which then
loads the WASM and makes a sample call to it with the output in console.

## What adding to the API looks like

In the sample, the API definition is [located here](proto/svc_sample.proto).

```protobuf
syntax = "proto3";

package sample;

message EchoRequest {
  string text = 1;
}

message EchoResponse {
  string text = 1;
}

// Define your WASM's API here
service WasmService {
  rpc Echo (EchoRequest) returns (EchoResponse) {}
}
```

Adding to the API might look like this:

```protobuf
syntax = "proto3";

package sample;

message EchoRequest {
  string text = 1;
}

message EchoResponse {
  string text = 1;
}

message AdditionRequest {
  int x = 1;
  int y = 2;
}

message AdditionResponse {
  int sum = 1;
}

// Define your WASM's API here
service WasmService {
  rpc Echo (EchoRequest) returns (EchoResponse) {}
  rpc Add (AdditionRequest) returns (AdditionResponse) {}
}
```

Running `make` will generate all the proto files for Go and Typescript.

Now you'd [add an implementation in your Go WASM service](lib/wasm/svc_wasm.go):

```go
func (s *wasmServer) Add(ctx context.Context, req *sample.AdditionRequest) (*sample.AdditionResponse, error) {
	return &sample.AdditionResponse {
		Sum: req.X + req.Y,
	}, nil
}
```

And you can call it [from Typescript](front/src/index.ts):

```typescript
  // After creating wasmService...
  const additionRequest = sample.AdditionRequest.create({ x: 3, y: 7 });
  const additionResponse = await wasmService.add(additionRequest);

  // Outputs '10' to the console
  console.log(additionResponse.sum);
```

Done!  Notice that any changes to the proto will create compile time errors
if either the front or back end isn't updated as well, giving delicious type
safety automatically.


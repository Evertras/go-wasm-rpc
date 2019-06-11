# go-wasm-rpc

Use gRPC service definitions to communicate between WASM and JS.
The intent is to allow you to use an existing toolset (gRPC, protobuf)
to allow a simple, type-safe API between JS (with Typescript) and
a Go WASM module.

With this setup, calls to WASM work just like a regular gRPC call using
[Protobufjs](https://github.com/protobufjs/protobuf.js).  The request is
serialized, passed to the WASM, and deserialized within.  The included
tool generates code to act as a translation layer for all this to happen
and should be part of the build pipeline.

Performance is decent but not great.  The overhead of memory copying is
not insignificant for high performance applications.  Fixing this is a
future goal that relies on a future release of Go.  See below for details.
The main benefit currently is strong type safety and ease of development.

This repo consists both of the tool itself and a sample for reference.
Run `make` to build everything.  Run `./sample-server` after `make` to
run a server that will serve the front end on `localhost:8000`, which then
loads the WASM and makes a sample call to it with the output in console.

## What using this API looks like

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
  int32 x = 1;
  int32 y = 2;
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

## Setting it up in a new project

There are a few things you need to do to use this in your own project.  This
assumes your project already has protobufs set up similar to this sample.
See [the Makefile](Makefile) for reference if needed.

First, add [the RPC implementation to your front end](front/src/rpc/impl.ts).
This can be simply copy/pasted with no changes.

Second, add the [service creation as in the sample shown here](front/src/index.ts).
Fit this where it makes sense for your project.  You only have to instantiate a single
service and can reuse it all you want.

Third, add [a service implementation](lib/wasm/svc_wasm.go) similar to this one
that instead matches your defined service from your proto.

Finally, generate the translation layer in Go by using the following:

```bash
# If you haven't run this yet
go get github.com/Evertras/go-wasm-rpc/cmd/wasm-rpc-gen

# Replace with your files/packages
wasm-rpc-gen lib/sample/svc_sample.pb.go lib/wasm/generated.go github.com/Evertras/go-wasm-rpc/lib/sample wasm
``` 

## Todo

As stated in the generated code, memory copying is currently awful.  A [future release of Go](https://github.com/golang/go/commit/c468ad04177c422534ad1ed4547295935f84743d)
will hopefully make this nicer.  There may be some tricks to do this without copying at all
by writing and reading directly to Go's memory, but that may be dangerous until the linked
commit goes live.

The actual usage of this tool is clunky.  It would be nicer to use [Cobra](https://github.com/spf13/cobra)
or similar down the line... maybe.  The pain is mitigated by the fact that it's a one time setup
that will be buried in a Makefile or equivalent somewhere.

The organization of this repo is a little messy.  Open to suggestions.

This README is a first pass and could almost certainly be improved and will need updating
when the organization of the repo is improved.


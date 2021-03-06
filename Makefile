TS_FILES=$(shell find front/src -name "*.ts")
WASM_FILES=$(shell find lib -name "*.go" ! -path "lib/static/*" ! -path "lib/server/*" ! -name "*_test.go")
GO_PROTO_BUILD_DIR=lib/sample
TS_PROTO_BUILD_DIR=front/src/proto

all: test build-sample-server build-sample-wasm front/index.js

clean:
	rm -rf $(GO_PROTO_BUILD_DIR)
	rm -rf $(TS_PROTO_BUILD_DIR)
	rm -f lib/wasm/generated.go

build-sample-server: protos
	CG_ENABLED=0 go build -o sample-server -v ./cmd/server/main.go

build-sample-wasm: protos front/lib.wasm
	GOARCH=wasm GOOS=js go build -o front/lib.wasm cmd/wasm/main.go

bench:
	go test -v -benchmem -bench . ./lib/...

run-dev:
	go run -race ./cmd/server/main.go -d

protos: $(GO_PROTO_BUILD_DIR) $(TS_PROTO_BUILD_DIR)

# These are not files, so always run them when asked to
.PHONY: all clean test build-sample-server build-sample-wasm bench run-dev protos

# Actual files/directories that must be generated
front/index.js: node_modules $(TS_PROTO_BUILD_DIR) $(TS_FILES)
	npx webpack || (rm -f front/index.js && exit 1)

node_modules:
	npm install

$(GO_PROTO_BUILD_DIR): proto/*.proto
	rm -rf $(GO_PROTO_BUILD_DIR)
	mkdir $(GO_PROTO_BUILD_DIR)
	@# Slightly weird PWD syntax here to deal with Windows gitbash mangling it otherwise.
	@# This is intentional, don't remove the initial slash!  You could also do this with
	@# protoc directly, but this is much more portable/consistent besides that quirk.
	docker run -v /${PWD}:/defs namely/protoc-all -d proto -l go -o $(GO_PROTO_BUILD_DIR) || (rm -rf $(GO_PROTO_BUILD_DIR) && exit 1)

$(TS_PROTO_BUILD_DIR): node_modules proto/*.proto
	rm -rf $(TS_PROTO_BUILD_DIR)
	mkdir $(TS_PROTO_BUILD_DIR)
	npx pbjs -t static-module -w commonjs proto/*.proto > $(TS_PROTO_BUILD_DIR)/proto.js || (rm -rf $(TS_PROTO_BUILD_DIR) && exit 1)
	npx pbts -o $(TS_PROTO_BUILD_DIR)/proto.d.ts $(TS_PROTO_BUILD_DIR)/proto.js || (rm -rf $(TS_PROTO_BUILD_DIR) && exit 1)

lib/wasm/generated.go: $(GO_PROTO_BUILD_DIR) cmd/wasm-rpc-gen/main.go
	go run ./cmd/wasm-rpc-gen/main.go -- lib/sample/svc_sample.pb.go lib/wasm/generated.go github.com/Evertras/go-wasm-rpc/lib/sample wasm

front/lib.wasm: $(WASM_FILES) cmd/wasm/main.go lib/wasm/generated.go
	GOARCH=wasm GOOS=js go build -o front/lib.wasm cmd/wasm/main.go


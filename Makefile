TS_FILES=$(shell find front/src -name "*.ts")
WASM_FILES=$(shell find lib -name "*.go" ! -path "lib/static/*" ! -path "lib/server/*" ! -name "*_test.go")
GO_PROTO_BUILD_DIR=lib/proto

all: test build

clean:
	rm -rf $(GO_PROTO_BUILD_DIR)

build-sample-server: protos
	CG_ENABLED=0 go build -o $(BINARY_NAME) -v ./cmd/server/main.go

build-sample-wasm: 
	GOARCH=wasm GOOS=js go build -o front/lib.wasm cmd/wasm/main.go

bench:
	go test -v -benchmem -bench . ./lib/...

run-dev:
	go run -race ./cmd/server/main.go -d

protos: $(GO_PROTO_BUILD_DIR) messages/tsmessage

# These are not files, so always run them when asked to
.PHONY: all clean test build-sample-server build-sample-wasm bench run-dev protos

# Actual files/directories that must be generated
front/index.js: node_modules messages/tsmessage $(TS_FILES)
	npx webpack || (rm -f front/index.js && exit 1)

node_modules:
	npm install

$(GO_PROTO_BUILD_DIR): messages/*.proto
	rm -rf $(GO_PROTO_BUILD_DIR)
	mkdir $(GO_PROTO_BUILD_DIR)
	@# Slightly weird PWD syntax here to deal with Windows gitbash mangling it otherwise.
	@# This is intentional, don't remove the initial slash!  You could also do this with
	@# protoc directly, but this is much more portable/consistent besides that quirk.
	docker run -v /${PWD}:/defs namely/protoc-all -d messages -l go -o $(GO_PROTO_BUILD_DIR) || (rm -rf $(GO_PROTO_BUILD_DIR) && exit 1)

front/src/proto: node_modules messages/*.proto
	rm -rf messages/tsmessage
	mkdir messages/tsmessage
	npx pbjs -t static-module -w commonjs messages/*.proto > messages/tsmessage/messages.js || (rm -rf messages/tsmessage && exit 1)
	npx pbts -o messages/tsmessage/messages.d.ts messages/tsmessage/messages.js || (rm -rf messages/tsmessage && exit 1)

front/lib.wasm: $(WASM_FILES) cmd/wasm/main.go
	GOARCH=wasm GOOS=js go build -o front/lib.wasm cmd/wasm/main.go


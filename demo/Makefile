.PHONY: copy

build: copy-wasm-exec
	GOOS=js GOARCH=wasm go build -o main.wasm main.go

copy-wasm-exec:
	cp "$(shell go env GOROOT)/lib/wasm/wasm_exec.js" .



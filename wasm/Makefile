build:
	GOOS=js GOARCH=wasm go build -o main.wasm

build-dev-web: build
	mv ./main.wasm ../app/todo-web/src/main.wasm
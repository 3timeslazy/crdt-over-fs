build:
	GOOS=js GOARCH=wasm go build -o main.wasm

build-dev-web: build
	mv ./main.wasm ../app/todo-web/src/main.wasm

run-minio:
	podman run \
		-p 9000:9000 \
		-p 9001:9001 \
		--name minio \
		-d \
		quay.io/minio/minio server /data --console-address ":9001"
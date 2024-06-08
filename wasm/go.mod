module github.com/3timeslazy/crdt-over-fs/wasm

go 1.22.3

replace github.com/3timeslazy/crdt-over-fs/sync => ../sync

require (
	github.com/3timeslazy/crdt-over-fs/sync v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go v1.53.19
)

require github.com/jmespath/go-jmespath v0.4.0 // indirect

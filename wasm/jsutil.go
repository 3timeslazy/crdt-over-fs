//go:build wasm

package main

import (
	"syscall/js"

	"github.com/3timeslazy/crdt-over-fs/sync"
)

func Promise(gofunc func() (any, error)) js.Value {
	asyncfn := js.FuncOf(func(this js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]

		go func() {
			v, err := gofunc()
			if err != nil {
				jserr := js.Global().Get("Error")
				reject.Invoke(jserr.New(err.Error()))
				return
			}

			if b, ok := v.(sync.State); ok {
				v = BytesToJS(b)
			}
			resolve.Invoke(v)
		}()

		return nil
	})

	p := js.Global().Get("Promise")
	return p.New(asyncfn)
}

func BytesFromJS(v js.Value) []byte {
	b := make([]byte, v.Length())
	js.CopyBytesToGo(b, v)
	return b
}

func BytesToJS(b []byte) js.Value {
	arr := js.Global().Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(arr, b)
	return arr
}

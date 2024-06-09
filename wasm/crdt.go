//go:build wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/3timeslazy/crdt-over-fs/sync"
)

type JSCRDT struct {
	jsObj js.Value
}

func (crdt *JSCRDT) EmptyState() sync.State {
	v := crdt.jsObj.Call("emptyState")
	if !v.InstanceOf(js.Global().Get("Uint8Array")) {
		panic(fmt.Sprintf(
			"emptyState must return %s, but got %s",
			"Uint8Array", v,
		))
	}

	return BytesFromJS(v)
}

func (crdt *JSCRDT) Merge(s1, s2 sync.State) (sync.State, []sync.Change, error) {
	v := crdt.jsObj.Call("merge", BytesToJS(s1), BytesToJS(s2))
	if v.Type() != js.TypeObject {
		panic(fmt.Sprintf(
			"merge must return %s, but got %s",
			js.TypeObject, v,
		))
	}

	// TODO: get changes
	state := BytesFromJS(v.Get("state"))

	return sync.State(state), nil, nil
}

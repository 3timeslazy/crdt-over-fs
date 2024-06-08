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
	if v.Type() != js.TypeString {
		panic(fmt.Sprintf(
			"emptyState must return %s, but got %s",
			js.TypeString, v,
		))
	}

	return sync.State(v.String())
}

func (crdt *JSCRDT) Merge(s1, s2 sync.State) (sync.State, []sync.Change, error) {
	v := crdt.jsObj.Call("merge", s1, s2)
	if v.Type() != js.TypeObject {
		panic(fmt.Sprintf(
			"merge must return %s, but got %s",
			js.TypeObject, v,
		))
	}

	// TODO: get changes
	state := v.Get("state").String()

	return sync.State(state), nil, nil
}

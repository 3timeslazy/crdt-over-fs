//go:build wasm

package main

import "io/fs"

// TODO: fs wrapper

type JSFS struct {
}

func (jsfs *JSFS) MakeDir(name string) error {
	panic("not implemented")
}

func (jsfs *JSFS) ReadDir(name string) ([]fs.DirEntry, error) {
	panic("not implemented")
}

func (jsfs *JSFS) WriteFile(name string, data []byte) error {
	panic("not implemented")
}

func (jsfs *JSFS) ReadFile(name string) ([]byte, error) {
	panic("not implemented")
}

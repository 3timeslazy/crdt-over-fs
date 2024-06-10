package fs

import "errors"

// TODO: add context?

var ErrNotExist = errors.New("file does not exist")

// FS is the interface that abstracts a
// file system.
//
// It's as simple and minimalist as possible to make
// it easier to implement that interfaces for as many
// possible storages as possible.
type FS interface {
	MakeDir(name string) error
	ReadDir(name string) ([]DirEntry, error)
	WriteFile(name string, data []byte) error
	ReadFile(name string) ([]byte, error)

	// TODO:
	// WatchDir(name, onChange)
}

type DirEntry interface {
	Name() string
	IsDir() bool
}

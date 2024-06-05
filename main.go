package main

import (
	"fmt"
	"io/fs"

	"github.com/automerge/automerge-go"
)

func main() {
	doc1 := automerge.New()

	err := doc1.RootMap().Set("hello", "world")
	if err != nil {
		panic(err)
	}

	doc2 := automerge.New()

	err = doc2.RootMap().Set("world", "hello")
	if err != nil {
		panic(err)
	}

	_, err = doc1.Merge(doc2)
	if err != nil {
		panic(err)
	}

	fmt.Println(doc1.RootMap().GoString())
}

type FS = fs.ReadDirFS

type Algorithm interface {
	Merge(state1, state2 []byte) ([]byte, error)
	// HumanReadable(state []byte) ([]byte, error)
}

type Automerge struct{}

func Merge(state1, state2 []byte) ([]byte, error) {
	d1, err := automerge.Load(state1)
	if err != nil {
		return nil, err
	}
	d2, err := automerge.Load(state2)
	if err != nil {
		return nil, err
	}
	_, err = d1.Merge(d2)
	if err != nil {
		return nil, err
	}
	return d1.Save(), nil
}

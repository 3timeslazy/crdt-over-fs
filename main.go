package main

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/automerge/automerge-go"
)

func main() {
	initial := automerge.New()
	initial.RootMap().Set("tasks", []Task{})

	s1 := state1(initial)
	s2 := state2(initial)

	Unwrap(s1.Merge(s2))

	merged := Unwrap(s1.RootMap().Get("tasks")).List()
	for _, v := range Unwrap(merged.Values()) {
		fmt.Println(v.GoString())
	}
	return

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

func state1(initial *automerge.Doc) *automerge.Doc {
	doc := Unwrap(initial.Fork())

	v := Unwrap(doc.RootMap().Get("tasks"))
	if v.IsVoid() {
		doc.RootMap().Set("tasks", []Task{})
	}

	tasks := Unwrap(doc.RootMap().Get("tasks")).List()

	Must(tasks.Append(Task{
		Name:      "Buy food",
		CreatedAt: time.Now(),
	}))

	return doc
}

func state2(initial *automerge.Doc) *automerge.Doc {
	doc := Unwrap(initial.Fork())

	v := Unwrap(doc.RootMap().Get("tasks"))
	if v.IsVoid() {
		doc.RootMap().Set("tasks", []Task{})
	}

	tasks := Unwrap(doc.RootMap().Get("tasks")).List()

	Must(tasks.Append(Task{
		Name:      "Buy PS5",
		CreatedAt: time.Now(),
	}))

	return doc
}

type Task struct {
	Name      string
	CreatedAt time.Time
}

func Unwrap[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

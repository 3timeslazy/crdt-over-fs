package main

import (
	"encoding/json"
	"fmt"

	"github.com/3timeslazy/crdt-over-fs/fs"
	"github.com/automerge/automerge-go"
)

type Tasks struct {
	doc  *automerge.Doc
	list *automerge.List
}

func NewTasks(initialState []byte) *Tasks {
	// TODO: return error

	var err error
	doc := automerge.New()

	if initialState != nil {
		doc, err = automerge.Load(initialState)
		if err != nil {
			panic(err)
		}
	}

	v, err := doc.RootMap().Get("tasks")
	if err != nil {
		panic(err)
	}
	if v.IsVoid() {
		doc.RootMap().Set("tasks", []Task{})
		_, err = doc.Commit("init tasks")
		if err != nil {
			panic(err)
		}
	}

	tasks, err := doc.RootMap().Get("tasks")
	if err != nil {
		panic(err)
	}

	return &Tasks{
		doc:  doc,
		list: tasks.List(),
	}
}

func (tasks Tasks) PushFront(task Task) {
	err := tasks.list.Insert(0, task)
	if err != nil {
		panic(err)
	}
	_, err = tasks.doc.Commit("added a task")
	if err != nil {
		panic(err)
	}
}

func (tasks Tasks) Remove(i int) {
	err := tasks.list.Delete(i)
	if err != nil {
		panic(err)
	}
	_, err = tasks.doc.Commit("removed a task")
	if err != nil {
		panic(err)
	}
}

func (tasks Tasks) All() []Task {
	values, err := tasks.list.Values()
	if err != nil {
		panic(err)
	}

	ts := []Task{}
	for _, v := range values {
		// TODO: I think it's possible to
		// bring an additional interface to the
		// automerge-go. Something similar to
		// json.Marshaler/Unmarshaler
		js, err := json.Marshal(v.Interface())
		if err != nil {
			panic(err)
		}
		t := Task{}
		err = json.Unmarshal(js, &t)
		if err != nil {
			panic(err)
		}
		ts = append(ts, t)
	}
	return ts
}

func (tasks Tasks) State() []byte {
	return tasks.doc.Save()
}

func (tasks Tasks) Merge(states []fs.State) {
	for _, state := range states {
		doc, err := automerge.Load(state)
		if err != nil {
			panic(err)
		}

		_, err = tasks.doc.Merge(doc)
		if err != nil {
			panic(err)
		}
	}
}

type Task struct {
	Name      string
	CreatedBy string
}

func (task Task) FilterValue() string {
	return task.Name
}

func (task Task) Title() string { return task.Name }

func (task Task) Description() string {
	s := fmt.Sprintf("by: %s\n", task.CreatedBy)
	return s
}

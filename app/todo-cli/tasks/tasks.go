package tasks

import (
	"errors"
	"fmt"

	"github.com/3timeslazy/crdt-over-fs/sync"

	"github.com/automerge/automerge-go"
)

type Manager struct {
	fs *sync.FSWrapper

	amdoc  *automerge.Doc
	amlist *automerge.List
}

// TODO: edit using the Text type?
// TODO: add counter
// TODO: save to remote on changes

type Task struct {
	Title  string `automerge:"title"`
	Author string `automerge:"author"`
	Done   bool   `automerge:"done"`
}

type Change struct {
	Message string
}

const (
	tasksKey = "tasks"

	authorKey = "author"
	titleKey  = "title"
	doneKey   = "done"
)

func NewManager(fs *sync.FSWrapper) *Manager {
	return &Manager{
		fs: fs,
	}
}

func (m *Manager) Init() error {
	state, err := m.fs.LoadOwnState()
	if err != nil {
		return err
	}

	amdoc, err := initState(state)
	if err != nil {
		return err
	}
	amlist := unwrap(amdoc.RootMap().Get(tasksKey)).List()

	m.amdoc = amdoc
	m.amlist = amlist

	return nil
}

func initState(state []byte) (*automerge.Doc, error) {
	doc, err := automerge.Load(state)
	if err != nil {
		return nil, err
	}
	v, err := doc.RootMap().Get(tasksKey)
	if err != nil {
		return nil, err
	}
	if v.IsVoid() {
		err = doc.RootMap().Set(tasksKey, []Task{})
		if err != nil {
			return nil, err
		}
		_, err = doc.Commit("init tasks list")
		if err != nil {
			return nil, err
		}
	}

	return doc, nil
}

func (m *Manager) Append(task Task) error {
	if task.Title == "" {
		return errors.New("title cannot be empty")
	}
	if task.Author == "" {
		return errors.New("author cannot be empty")
	}

	m.changeAndCommit(
		fmt.Sprintf("%q added task %q", task.Author, task.Title),
		func(list *automerge.List) error {
			return list.Append(task)
		},
	)
	return nil
}

func (m *Manager) PushFront(task Task) error {
	if task.Title == "" {
		return errors.New("title cannot be empty")
	}
	if task.Author == "" {
		return errors.New("author cannot be empty")
	}

	m.changeAndCommit(
		fmt.Sprintf("%q added task %q", task.Author, task.Title),
		func(list *automerge.List) error {
			return list.Insert(0, task)
		},
	)
	return nil
}

func (m *Manager) Remove(i int) error {
	if i < 0 || i >= m.amlist.Len() {
		return errors.New("index out of range")
	}

	task := unwrap(automerge.As[Task](
		m.amdoc.Path(tasksKey, i).Get(),
	))
	m.changeAndCommit(
		fmt.Sprintf("%q removed task %q", task.Author, task.Title),
		func(list *automerge.List) error {
			return list.Delete(i)
		},
	)
	return nil
}

func (m *Manager) SetTitle(i int, newTitle string) error {
	if i < 0 || i >= m.amlist.Len() {
		return errors.New("index out of range")
	}
	if newTitle == "" {
		return errors.New("new title cannot be empty")
	}

	task := unwrap(automerge.As[Task](
		m.amdoc.Path(tasksKey, i).Get(),
	))
	m.changeAndCommit(
		fmt.Sprintf("%q changed title from %q to %q", task.Author, task.Title, newTitle),
		func(list *automerge.List) error {
			task := unwrap(list.Get(i)).Map()
			return task.Set(titleKey, newTitle)
		},
	)
	return nil
}

func (m *Manager) ToggleDone(i int) error {
	if i < 0 || i >= m.amlist.Len() {
		return errors.New("index out of range")
	}

	task := unwrap(automerge.As[Task](
		m.amdoc.Path(tasksKey, i).Get(),
	))
	status := "done"
	if task.Done {
		status = "not done"
	}
	m.changeAndCommit(
		fmt.Sprintf("%q marked %q as %q", task.Author, task.Title, status),
		func(list *automerge.List) error {
			amtask := unwrap(list.Get(i)).Map()
			return amtask.Set(doneKey, !task.Done)
		},
	)
	return nil
}

func (m *Manager) Persist() error {
	return m.fs.SaveOwnState(m.amdoc.Save())
}

func (m *Manager) MergeRemote() ([]Change, error) {
	newState, allChanges, err := m.fs.Sync(m.amdoc.Save())
	if err != nil {
		return nil, err
	}

	newDoc, err := automerge.Load(newState)
	if err != nil {
		return nil, err
	}
	m.amdoc = newDoc
	m.amlist = unwrap(m.amdoc.RootMap().Get(tasksKey)).List()

	changes := []Change{}
	for _, neighbourChanges := range allChanges {
		for _, change := range neighbourChanges {
			// TODO: don't panic
			hash := unwrap(automerge.NewChangeHash(change.Hash))
			change := unwrap(m.amdoc.Change(hash))
			changes = append(changes, Change{
				Message: change.Message(),
			})
		}
	}

	return changes, nil
}

func (m *Manager) changeAndCommit(
	msg string,
	change func(list *automerge.List) error,
) {
	err := change(m.amlist)
	if err != nil {
		panic(err)
	}

	_, err = m.amdoc.Commit(msg)
	if err != nil {
		panic(err)
	}
}

func Map[T any](m *Manager, fn func(Task) T) []T {
	amValues := unwrap(m.amlist.Values())
	out := make([]T, 0, len(amValues))

	for _, amv := range amValues {
		v := unwrap(automerge.As[Task](amv))
		out = append(out, fn(v))
	}

	return out
}

func GetAs[T any](m *Manager, i int, fn func(Task) T) T {
	amv := unwrap(m.amlist.Get(i))
	task := unwrap(automerge.As[Task](amv))
	return fn(task)
}

func unwrap[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

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
// TODO: add done
// TODO: save to remote on changes

type Task struct {
	Title  string `automerge:"title"`
	Author string `automerge:"author"`
	Done   bool   `automerge:"-"`
}

type Change struct {
	Message string
}

const (
	tasksKey  = "tasks"
	authorKey = "author"
	titleKey  = "title"
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

	author := unwrap(automerge.As[string](
		m.amdoc.Path(tasksKey, i, authorKey).Get(),
	))
	title := unwrap(automerge.As[string](
		m.amdoc.Path(tasksKey, i, titleKey).Get(),
	))

	m.changeAndCommit(
		fmt.Sprintf("%q removed task %q", author, title),
		func(list *automerge.List) error {
			return list.Delete(i)
		},
	)
	return nil
}

func (m *Manager) SetTitle(taskIdx int, newTitle string) error {
	if taskIdx < 0 || taskIdx >= m.amlist.Len() {
		return errors.New("index out of range")
	}
	if newTitle == "" {
		return errors.New("new title cannot be empty")
	}

	author := unwrap(automerge.As[string](
		m.amdoc.Path(tasksKey, taskIdx, authorKey).Get(),
	))
	title := unwrap(automerge.As[string](
		m.amdoc.Path(tasksKey, taskIdx, titleKey).Get(),
	))

	m.changeAndCommit(
		fmt.Sprintf("%q changed title from %q to %q", author, title, newTitle),
		func(list *automerge.List) error {
			task := unwrap(list.Get(taskIdx)).Map()
			return task.Set(titleKey, newTitle)
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
	change func(doc *automerge.List) error,
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

func unwrap[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

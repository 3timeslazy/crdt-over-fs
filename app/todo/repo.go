package main

import (
	"github.com/3timeslazy/crdt-over-fs/sync"
)

type Repository struct {
	fs *sync.FSWrapper
}

func NewRepository(fs *sync.FSWrapper) *Repository {
	return &Repository{
		fs: fs,
	}
}

func (repo *Repository) LoadTasks() (*Tasks, error) {
	state, err := repo.fs.LoadOwnState()
	if err != nil {
		return nil, err
	}

	return TasksFromState(state), nil
}

func (wrapper *Repository) SaveTasks(tasks *Tasks) error {
	return wrapper.fs.SaveOwnState(tasks.State())
}

func (repo *Repository) Sync(tasks *Tasks) (*Tasks, map[string][]sync.Change, error) {
	newState, changes, err := repo.fs.Sync(tasks.State())
	if err != nil {
		return nil, nil, err
	}

	return TasksFromState(newState), changes, nil
}

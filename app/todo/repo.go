package main

import (
	"errors"
	"fmt"

	"github.com/3timeslazy/crdt-over-fs/sync"
	"github.com/automerge/automerge-go"
)

type Repository struct {
	id string
	fs *sync.FSWrapper
}

func NewRepository(stateID string, fs *sync.FSWrapper) *Repository {
	return &Repository{
		id: stateID,
		fs: fs,
	}
}

func (repo *Repository) LoadTasks() (*Tasks, error) {
	state, err := repo.fs.LoadOwnState()
	if err != nil && !errors.Is(err, sync.ErrStateNotFound) {
		return nil, fmt.Errorf("load state: %w", err)
	}
	if err == nil {
		return NewTasks(state), nil
	}

	// ErrStateNotFound case
	//
	// This case is important, because if we just create
	// a new document that will be the same as removing everything

	neighbours, _, err := repo.fs.LoadNeighbourStates()
	if err != nil {
		return nil, fmt.Errorf("load neighbour states: %w", err)
	}
	if len(neighbours) == 0 {
		return NewTasks(nil), nil
	}

	initial, err := automerge.Load(neighbours[0])
	if err != nil {
		return nil, fmt.Errorf("load 0th neighbour state: %w", err)
	}

	for _, state := range neighbours[1:] {
		doc, err := automerge.Load(state)
		if err != nil {
			return nil, fmt.Errorf("load neighbour state: %w", err)
		}
		_, err = initial.Merge(doc)
		if err != nil {
			return nil, fmt.Errorf("merge states: %w", err)
		}
	}

	return NewTasks(initial.Save()), nil
}

func (wrapper *Repository) SaveTasks(tasks *Tasks) error {
	return wrapper.fs.SaveOwnState(tasks.State())
}

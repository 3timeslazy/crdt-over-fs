package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"

	stdfs "io/fs"

	"github.com/3timeslazy/crdt-over-fs/fs"
	"github.com/3timeslazy/crdt-over-fs/fs/local"
)

type FSWrapper struct {
	id string
	fs fs.FS
}

const appDir = "./todo-over-fs"

func NewFSWrapper(id string) *FSWrapper {
	return &FSWrapper{
		id: id,
		fs: &local.FS{},
	}
}

func (wrapper *FSWrapper) SetupDir() error {
	_, err := wrapper.fs.ReadDir(appDir)
	if errors.Is(err, stdfs.ErrNotExist) {
		return wrapper.fs.MakeDir(appDir)
	}
	return err
}

func (wrapper *FSWrapper) LoadTasks() ([]Task, error) {
	entries, err := wrapper.fs.ReadDir(appDir)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() != wrapper.id {
			continue
		}

		filepath := path.Join(appDir, entry.Name())

		state, err := wrapper.fs.ReadFile(filepath)
		if err != nil {
			return nil, err
		}

		tasks := []Task{}
		err = json.Unmarshal(state, &tasks)
		if err != nil {
			return nil, err
		}

		return tasks, nil
	}

	return []Task{}, nil
}

func (wrapper *FSWrapper) SaveTasks(tasks []Task) error {
	state, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	stateFilepath := path.Join(appDir, wrapper.id)
	return wrapper.fs.WriteFile(stateFilepath, state)
}

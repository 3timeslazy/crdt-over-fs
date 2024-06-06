package fs

import (
	"errors"
	"fmt"
	iofs "io/fs"
	"path"
)

type Wrapper struct {
	StateID string
	RootDir string
	FS      FS
}

type State []byte

func (w *Wrapper) SetupDir() error {
	_, err := w.FS.ReadDir(w.RootDir)
	// TODO: replace iofs.ErrNotExist with a ErrNotFound custom error.
	if errors.Is(err, iofs.ErrNotExist) {
		return w.FS.MakeDir(w.RootDir)
	}
	return err
}

var ErrStateNotFound = errors.New("no state found")

func (w *Wrapper) LoadOwnState() (State, error) {
	entries, err := w.FS.ReadDir(w.RootDir)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() != w.StateID {
			continue
		}

		filepath := path.Join(w.RootDir, entry.Name())
		state, err := w.FS.ReadFile(filepath)
		if err != nil {
			return nil, err
		}

		return state, nil
	}

	return nil, ErrStateNotFound
}

func (w *Wrapper) SaveOwnState(state State) error {
	stateFilepath := path.Join(w.RootDir, w.StateID)
	return w.FS.WriteFile(stateFilepath, state)
}

func (w *Wrapper) LoadNeighbourStates() ([]State, []string, error) {
	entries, err := w.FS.ReadDir(w.RootDir)
	if err != nil {
		return nil, nil, err
	}

	neighbours := []State{}
	ids := []string{}
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == w.StateID {
			continue
		}

		filepath := path.Join(w.RootDir, entry.Name())
		state, err := w.FS.ReadFile(filepath)
		if err != nil {
			return nil, nil, err
		}

		neighbours = append(neighbours, state)
		ids = append(ids, entry.Name())
	}

	return neighbours, ids, nil
}

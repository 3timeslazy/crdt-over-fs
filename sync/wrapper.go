package sync

import (
	"errors"
	"fmt"
	"path"

	"github.com/3timeslazy/crdt-over-fs/sync/fs"
)

// FSWrapper provides high-level state operations
// using underlying file system.
//
// Its goal is to be a common solution for any app
// using synchronisation via combination of
// CRDT and file system.
type FSWrapper struct {
	fs      fs.FS
	crdt    CRDT
	stateID string
	rootDir string
}

// TODO: add work with local persistent layer such as local files or
// localStorage on frontend.

func NewFSWrapper(fs fs.FS, crdt CRDT, stateID, rootDir string) *FSWrapper {
	return &FSWrapper{
		fs:      fs,
		crdt:    crdt,
		stateID: stateID,
		rootDir: rootDir,
	}
}

func (w *FSWrapper) InitRootDir() error {
	_, err := w.fs.ReadDir(w.rootDir)
	if errors.Is(err, fs.ErrNotExist) {
		return w.fs.MakeDir(w.rootDir)
	}
	return err
}

func (w *FSWrapper) LoadOwnState() (State, error) {
	filepath := path.Join(w.rootDir, w.stateID)
	state, err := w.fs.ReadFile(filepath)
	if err == nil {
		return state, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("read state file: %w", err)
	}

	// We come here when there is no own state found. This case
	// is important, because if we just create a new document
	// that will be the same as removing everything
	neighbours, _, err := w.loadNeighbourStates()
	if err != nil {
		return nil, fmt.Errorf("load neighbour states: %w", err)
	}
	if len(neighbours) == 0 {
		return w.crdt.EmptyState(), nil
	}

	initial := neighbours[0]
	for _, state := range neighbours[1:] {
		merged, _, err := w.crdt.Merge(initial, state)
		if err != nil {
			return nil, fmt.Errorf("merge neighbour state: %w", err)
		}
		initial = merged
	}

	return initial, nil
}

func (w *FSWrapper) SaveOwnState(state State) error {
	stateFilepath := path.Join(w.rootDir, w.stateID)
	return w.fs.WriteFile(stateFilepath, state)
}

func (w *FSWrapper) loadNeighbourStates() ([]State, []string, error) {
	entries, err := w.fs.ReadDir(w.rootDir)
	if err != nil {
		return nil, nil, err
	}

	neighbours := []State{}
	ids := []string{}
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == w.stateID {
			continue
		}

		filepath := path.Join(w.rootDir, entry.Name())
		state, err := w.fs.ReadFile(filepath)
		if err != nil {
			return nil, nil, err
		}

		neighbours = append(neighbours, state)
		ids = append(ids, entry.Name())
	}

	return neighbours, ids, nil
}

func (w *FSWrapper) Sync(localState State) (State, map[string][]Change, error) {
	neighbours, ids, err := w.loadNeighbourStates()
	if err != nil {
		return nil, nil, err
	}

	totalChanges := map[string][]Change{}
	for i, neighbour := range neighbours {
		state, changes, err := w.crdt.Merge(localState, neighbour)
		if err != nil {
			return nil, nil, err
		}

		totalChanges[ids[i]] = changes
		localState = state
	}

	return localState, totalChanges, nil
}

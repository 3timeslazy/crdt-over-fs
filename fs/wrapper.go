package fs

import (
	"errors"
	"fmt"
	"path"
)

// Wrapper provides high-level state operations
// using underlying file system.
//
// Its goal is to be a common solution for any app
// using synchronisation via combination of
// CRDT and file system.
type Wrapper struct {
	fs      FS
	stateID string
	rootDir string
}

// State is a CRDT representations of the app's state.
type State []byte

func NewWrapper(fs FS, stateID, rootDir string) *Wrapper {
	return &Wrapper{
		fs:      fs,
		stateID: stateID,
		rootDir: rootDir,
	}
}

func (w *Wrapper) InitRootDir() error {
	_, err := w.fs.ReadDir(w.rootDir)
	if errors.Is(err, ErrNotExist) {
		return w.fs.MakeDir(w.rootDir)
	}
	return err
}

var ErrStateNotFound = errors.New("no state found")

func (w *Wrapper) LoadOwnState() (State, error) {
	filepath := path.Join(w.rootDir, w.stateID)
	state, err := w.fs.ReadFile(filepath)
	if errors.Is(err, ErrNotExist) {
		return nil, ErrStateNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("read state file: %w", err)
	}

	return state, nil
}

func (w *Wrapper) SaveOwnState(state State) error {
	stateFilepath := path.Join(w.rootDir, w.stateID)
	return w.fs.WriteFile(stateFilepath, state)
}

func (w *Wrapper) LoadNeighbourStates() ([]State, []string, error) {
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

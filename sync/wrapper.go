package sync

import (
	"errors"
	"fmt"
	"path"

	"github.com/3timeslazy/crdt-over-fs/fs"
)

var (
	ErrStateNotFound = errors.New("no state found")
)

// FSWrapper provides high-level state operations
// using underlying file system.
//
// Its goal is to be a common solution for any app
// using synchronisation via combination of
// CRDT and file system.
type FSWrapper struct {
	fs      fs.FS
	stateID string
	rootDir string
}

// State is a CRDT representations of the app's state.
type State []byte

func NewFSWrapper(fs fs.FS, stateID, rootDir string) *FSWrapper {
	return &FSWrapper{
		fs:      fs,
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
	if errors.Is(err, fs.ErrNotExist) {
		return nil, ErrStateNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("read state file: %w", err)
	}

	return state, nil
}

func (w *FSWrapper) SaveOwnState(state State) error {
	stateFilepath := path.Join(w.rootDir, w.stateID)
	return w.fs.WriteFile(stateFilepath, state)
}

func (w *FSWrapper) LoadNeighbourStates() ([]State, []string, error) {
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

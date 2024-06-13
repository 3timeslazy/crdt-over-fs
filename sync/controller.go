package sync

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/3timeslazy/crdt-over-fs/sync/fs"
)

// TODO: delete strategy

// TODO: additional functions for work with assets such as pictures or video. Since
// the code uses a file system it might store assets there as well.

// TODO: local persistent layer such as local files or localStorage on frontend?

// TODO: better filename format. Examples:
//   - ${rootDir}/states/...

// TODO: tests

// Controller provides high-level state operations
// using underlying file system.
//
// Its goal is to be a common solution for any app
// using synchronisation via combination of
// CRDT and file system.
type Controller struct {
	fs      fs.FS
	crdt    CRDT
	stateID string
	rootDir string
}

type StateFile struct {
	State        []byte    `json:"state"`
	LastModified time.Time `json:"lastModified"`
}

func NewController(fs fs.FS, crdt CRDT, stateID, rootDir string) *Controller {
	return &Controller{
		fs:      fs,
		crdt:    crdt,
		stateID: stateID + ".json",
		rootDir: rootDir,
	}
}

func (ctrl *Controller) InitRootDir() error {
	_, err := ctrl.fs.ReadDir(ctrl.rootDir)
	if errors.Is(err, fs.ErrNotExist) {
		return ctrl.fs.MakeDir(ctrl.rootDir)
	}
	return err
}

func (ctrl *Controller) LoadOwnState() (State, error) {
	filepath := path.Join(ctrl.rootDir, ctrl.stateID)
	fileContent, err := ctrl.fs.ReadFile(filepath)
	if err == nil {
		file := StateFile{}
		err := json.Unmarshal(fileContent, &file)
		if err != nil {
			return nil, fmt.Errorf("parse json: %w", err)
		}
		return file.State, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("read state file: %w", err)
	}

	// We come here when there is no own state found. This case
	// is important, because if we just create a new document
	// that will be the same as removing everything
	neighbours, _, err := ctrl.loadNeighbourFiles()
	if err != nil {
		return nil, fmt.Errorf("load neighbour states: %w", err)
	}
	if len(neighbours) == 0 {
		return ctrl.crdt.EmptyState(), nil
	}

	initial := neighbours[0].State
	for _, state := range neighbours[1:] {
		merged, _, err := ctrl.crdt.Merge(initial, state.State)
		if err != nil {
			return nil, fmt.Errorf("merge neighbour state: %w", err)
		}
		initial = merged
	}

	return initial, nil
}

func (ctrl *Controller) SaveOwnState(state State) error {
	file := StateFile{
		State:        state,
		LastModified: time.Now(),
	}
	content, err := json.Marshal(file)
	if err != nil {
		return err
	}

	stateFilepath := path.Join(ctrl.rootDir, ctrl.stateID)
	return ctrl.fs.WriteFile(stateFilepath, content)
}

func (ctrl *Controller) loadNeighbourFiles() ([]StateFile, []string, error) {
	entries, err := ctrl.fs.ReadDir(ctrl.rootDir)
	if err != nil {
		return nil, nil, err
	}

	neighbours := []StateFile{}
	ids := []string{}
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == ctrl.stateID {
			continue
		}

		filepath := path.Join(ctrl.rootDir, entry.Name())
		state, err := ctrl.fs.ReadFile(filepath)
		if err != nil {
			return nil, nil, err
		}
		file := StateFile{}
		err = json.Unmarshal(state, &file)
		if err != nil {
			return nil, nil, fmt.Errorf("parse json: %w", err)
		}

		neighbours = append(neighbours, file)
		ids = append(ids, entry.Name())
	}

	return neighbours, ids, nil
}

func (ctrl *Controller) Sync(localState State) (State, map[string][]Change, error) {
	neighbours, ids, err := ctrl.loadNeighbourFiles()
	if err != nil {
		return nil, nil, err
	}

	totalChanges := map[string][]Change{}
	for i, neighbour := range neighbours {
		merged, changes, err := ctrl.crdt.Merge(localState, neighbour.State)
		if err != nil {
			return nil, nil, err
		}

		totalChanges[ids[i]] = changes
		localState = merged
	}

	return localState, totalChanges, nil
}

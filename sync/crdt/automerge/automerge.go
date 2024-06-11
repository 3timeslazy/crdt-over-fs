package automerge

import (
	"fmt"

	"github.com/3timeslazy/crdt-over-fs/sync"
	"github.com/automerge/automerge-go"
)

type Automerge struct{}

func (am *Automerge) EmptyState() sync.State {
	return automerge.New().Save()
}

func (am *Automerge) Merge(local, remote sync.State) (sync.State, []sync.Change, error) {
	localDoc, err := automerge.Load(local)
	if err != nil {
		return nil, nil, fmt.Errorf("load local state: %w", err)
	}
	heads := localDoc.Heads()

	remoteDoc, err := automerge.Load(remote)
	if err != nil {
		return nil, nil, fmt.Errorf("load remote state: %w", err)
	}
	// Do not use changes returned by Merge, because it
	// returns changes even if the remote state has already
	// been merged into the local one.
	//
	// An example:
	//   1. local, remote = initialState
	//   2. change(remote)
	//   3. local = merge(local, remote)
	//       --> a change from 2. returned
	//   4. local = merge(local, remote)
	//       --> the same change from 3. returned again
	//
	_, err = localDoc.Merge(remoteDoc)
	if err != nil {
		return nil, nil, fmt.Errorf("merge remote state: %w", err)
	}

	changes, err := localDoc.Changes(heads...)
	if err != nil {
		return nil, nil, fmt.Errorf("retrieve new changes: %w", err)
	}
	syncChanges := make([]sync.Change, 0, len(changes))
	for _, c := range changes {
		syncChanges = append(syncChanges, sync.Change{
			Hash: c.Hash().String(),
		})
	}
	return localDoc.Save(), syncChanges, nil
}

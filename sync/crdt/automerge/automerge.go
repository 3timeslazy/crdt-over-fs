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

func (am *Automerge) Merge(s1, s2 sync.State) (sync.State, []sync.Change, error) {
	d1, err := automerge.Load(s1)
	if err != nil {
		return nil, nil, fmt.Errorf("load state 1: %w", err)
	}
	d2, err := automerge.Load(s2)
	if err != nil {
		return nil, nil, fmt.Errorf("load state 2: %w", err)
	}
	changes, err := d1.Merge(d2)
	if err != nil {
		return nil, nil, fmt.Errorf("merge states: %w", err)
	}
	syncChanges := make([]sync.Change, 0, len(changes))
	for _, c := range changes {
		syncChanges = append(syncChanges, sync.Change{
			Hash: c.String(),
		})
	}
	return d1.Save(), syncChanges, nil
}

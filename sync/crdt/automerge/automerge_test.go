package automerge

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/automerge/automerge-go"
)

func TestMerge_MergeChange(t *testing.T) {
	crdt := &Automerge{}
	initialState := crdt.EmptyState()

	remote := unwrap(automerge.Load(initialState))
	local := unwrap(automerge.Load(initialState))

	assert.NoError(t, remote.RootMap().Set("hello", "world"))
	_, err := remote.Commit("hello world")
	assert.NoError(t, err)

	merged, syncChanges, err := crdt.Merge(local.Save(), remote.Save())
	assert.NoError(t, err)
	local, err = automerge.Load(merged)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(syncChanges))
	hash, err := automerge.NewChangeHash(syncChanges[0].Hash)
	assert.NoError(t, err)

	change, err := local.Change(hash)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", change.Message())
}

func TestMerge_MergeSameRemoteTwice(t *testing.T) {
	crdt := &Automerge{}
	initialState := crdt.EmptyState()

	remote := unwrap(automerge.Load(initialState))
	local := unwrap(automerge.Load(initialState))

	assert.NoError(t, remote.RootMap().Set("hello", "world"))
	_, err := remote.Commit("hello world")
	assert.NoError(t, err)

	merged, syncChanges, err := crdt.Merge(local.Save(), remote.Save())
	assert.NoError(t, err)
	local, err = automerge.Load(merged)
	assert.NoError(t, err)

	// Since it's the first merge of the remote into the local
	// a change must be returned.
	assert.Equal(t, 1, len(syncChanges))

	hash, err := automerge.NewChangeHash(syncChanges[0].Hash)
	assert.NoError(t, err)
	change, err := local.Change(hash)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", change.Message())

	// Call `Merge` on the the remote twice should not
	// return new changes
	_, syncChanges, err = crdt.Merge(local.Save(), remote.Save())
	assert.NoError(t, err)
	assert.Equal(t, 0, len(syncChanges))
}

func unwrap[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

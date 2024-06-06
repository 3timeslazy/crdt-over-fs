package sync

type CRDT interface {
	EmptyState() State
	Merge(s1, s2 State) (State, []Change, error)
}

// State is a raw CRDT representation of an app's state.
type State []byte

// TODO: More informative changes

type Change struct {
	Hash string
}

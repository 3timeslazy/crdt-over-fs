package fs

// TODO: add context?

type FS interface {
	MakeDir(name string) error
	ReadDir(name string) ([]DirEntry, error)
	WriteFile(name string, data []byte) error
	ReadFile(name string) ([]byte, error)
}

type DirEntry interface {
	Name() string
	IsDir() bool
}

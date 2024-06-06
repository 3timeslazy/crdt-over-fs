package local

import (
	"os"

	"github.com/3timeslazy/crdt-over-fs/fs"
)

type FS struct{}

func NewFS() *FS {
	return &FS{}
}

func (localfs *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	osEntries, err := os.ReadDir(name)
	if err != nil {
		return nil, err
	}

	entries := []fs.DirEntry{}
	for _, entry := range osEntries {
		entries = append(entries, entry)
	}
	return entries, nil
}

func (localfs *FS) MakeDir(name string) error {
	return os.Mkdir(name, 0777)
}

func (localfs *FS) WriteFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0666)
}

func (localfs *FS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

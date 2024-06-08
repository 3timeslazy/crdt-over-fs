package local

import (
	"errors"
	"os"

	"github.com/3timeslazy/crdt-over-fs/sync/fs"
)

type FS struct{}

func NewFS() *FS {
	return &FS{}
}

func (localfs *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	osEntries, err := os.ReadDir(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fs.ErrNotExist
		}
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
	file, err := os.ReadFile(name)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fs.ErrNotExist
	}
	return file, err
}

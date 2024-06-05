package local

import (
	"os"

	"github.com/3timeslazy/crdt-over-fs/fs"
)

type FS struct{}

func (localfs *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
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

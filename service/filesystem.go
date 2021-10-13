package service

import (
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

// FileSystem custom file system handler
type FileSystem struct {
	fs http.FileSystem
}

// Open opens file
func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "fs.Open")
	}

	s, err := f.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "f.Stat")
	}
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			return nil, errors.Wrap(err, "fs.Open")
		}
	}

	return f, nil
}

package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrNoMoreFiles = errors.New("no more files in the walker")

type DirWalker interface {
	Next() (*FileEntry, error)
}

type FileEntry struct {
	Path     []string // relative to origin
	Language string
	FullPath string
}

func (fe FileEntry) ReadContents() ([]byte, error) {
	return os.ReadFile(fe.FullPath)
}

type ioDirWalker struct {
	Origin  string
	files   []FileEntry
	current int
}

func IoDirWalker(dir string) (*ioDirWalker, error) {
	walker := &ioDirWalker{Origin: dir, current: -1}
	err := walker.loadFiles()
	if err != nil {
		return nil, err
	}
	return walker, nil
}

func (walker *ioDirWalker) loadFiles() error {
	err := filepath.WalkDir(walker.Origin, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".json") {
			relPath, err := filepath.Rel(walker.Origin, path)
			if err != nil {
				return err
			}

			// Get the path components and remove the extension from the name
			dirPath := filepath.Dir(relPath)
			if dirPath == "." {
				dirPath = ""
			}
			// Use empty slice if path is "."
			relativePath := filepath.SplitList(dirPath)
			if len(relativePath) == 1 && relativePath[0] == "" {
				relativePath = []string{}
			}

			// Save only the file name without the extension
			fileNameWithoutExt := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))

			walker.files = append(walker.files, FileEntry{
				Path:     relativePath,
				Language: fileNameWithoutExt,
				FullPath: path,
			})
		}
		return nil
	})
	return err
}

func (walker *ioDirWalker) Next() (*FileEntry, error) {
	walker.current++
	if walker.current >= len(walker.files) || walker.files == nil {
		return nil, ErrNoMoreFiles
	}
	return &walker.files[walker.current], nil
}

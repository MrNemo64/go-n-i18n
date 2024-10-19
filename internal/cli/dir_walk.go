package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrNoMoreFiles = errors.New("no more files in the walker")

type DirWalker interface {
	Next() (FileEntry, error)
}

type FileEntry interface {
	Path() []string
	Language() string
	FullPath() string
	ReadContents() ([]byte, error)
}

type IOFileEntry struct {
	path     []string
	language string
	fullPath string
}

func (fe *IOFileEntry) Path() []string                { return fe.path }
func (fe *IOFileEntry) Language() string              { return fe.language }
func (fe *IOFileEntry) FullPath() string              { return fe.fullPath }
func (fe *IOFileEntry) ReadContents() ([]byte, error) { return os.ReadFile(fe.fullPath) }

type ioDirWalker struct {
	Origin  string
	files   []IOFileEntry
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

			walker.files = append(walker.files, IOFileEntry{
				path:     relativePath,
				language: fileNameWithoutExt,
				fullPath: path,
			})
		}
		return nil
	})
	return err
}

func (walker *ioDirWalker) Next() (FileEntry, error) {
	walker.current++
	if walker.current >= len(walker.files) || walker.files == nil {
		return nil, ErrNoMoreFiles
	}
	return &walker.files[walker.current], nil
}

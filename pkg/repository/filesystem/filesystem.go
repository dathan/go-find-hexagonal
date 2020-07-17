package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dathan/go-find-hexagonal/pkg/find"
)

// FileSystemRepository is a struct to the backend which performs the find - the filesystem
type FilesytemRepository struct {
}

// NewFileSystemRespository  returns the struct
func NewFileSystemRepository() *FilesytemRepository {
	ret := &FilesytemRepository{}
	return ret
}

// Implements the repository interface
func (f *FilesytemRepository) Find(fo find.FilterOptions) (find.FindResults, error) {
	path := fo.GetStart()

	findResults := find.FindResults{}

	// todo: investigate if I need to lock the struct
	walkFn := func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}
		if path == info.Name() {
			path = "."
		}

		fResult := find.FindResult{
			Name:      info.Name(),
			CreatedAt: info.ModTime().Unix(),
			Path:      path,
		}

		fFunc := fo.GetFilterFunc()
		if fFunc != nil && fFunc(&fResult) {
			findResults = append(findResults, fResult)
		} else {
			fmt.Printf("Skipping this results :%v\n", fResult)
		}

		return nil
	}

	if err := filepath.Walk(path, walkFn); err != nil {
		fmt.Printf("ERROR: filepath.Walk: %s\n", err)
		return nil, err
	}

	return findResults, nil
}

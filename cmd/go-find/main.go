//
//
// find.New().Positional().Global().Tests().Operators().Do()
// or
// find.New().Tests().Actions().Operators()
// or
// Tests(): //-type -atime,ctime,mtime,-name,-path,-perm,-size,-user,-uid,
// Actions(): // -delete,exec
// Operators(): // expr -o expr ||
// find.New().Type(<string>).Time(<TIMEIMPL>).Name('<STRING>').Do()
package main

import (
	"fmt"

	"github.com/dathan/go-find-hexagonal/pkg/filter"
	"github.com/dathan/go-find-hexagonal/pkg/find"
	"github.com/dathan/go-find-hexagonal/pkg/repository/filesystem"
)

func main() {

	// let compiler know know that your using an interface
	var fileSystemFilter find.FilterOptions
	var repository find.Repository
	var fAbj find.Find

	// ensure that FilterFunc filters out the things your looking for
	var fFunc find.FilterFunc = func(fr *find.FindResult) bool {
		if fr != nil && fr.Name == "filesystem.go" {
			return true
		}
		return false
	}

	// set the filter func to find the specific files
	fileSystemFilter = filter.NewGenericOptions(".").SetFilterFunc(fFunc)

	// this will not compile if the package did not implement the inteface.
	repository = filesystem.NewFileSystemRepository()

	fAbj = find.New(repository)

	fr, err := fAbj.Do(fileSystemFilter)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
	fmt.Println(fr)

}

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
	"flag"
	"fmt"
	"strings"

	"github.com/dathan/go-find-hexagonal/pkg/filter"
	"github.com/dathan/go-find-hexagonal/pkg/find"
	"github.com/dathan/go-find-hexagonal/pkg/repository/filesystem"
	social "github.com/dathan/go-find-hexagonal/pkg/repository/social"
	"github.com/dathan/go-find-hexagonal/pkg/ui"
)

func main() {

	// flags
	var root string
	var name string
	var findType string

	flag.StringVar(&root, "root", "dathanvp", "depending on the type search filesystem names or twitter handle favorites")
	flag.StringVar(&name, "name", "", "giving a name look for that from the root")
	flag.StringVar(&findType, "type", "social", "Social Bookmarks lookup or filesystem finds")
	flag.Parse()

	// let compiler know know that your using an interface
	var filterOpt find.FilterOptions
	var repository find.Repository
	var fAbj find.Find

	// ensure that FilterFunc filters out the things your looking for
	var fFunc find.FilterFunc = func(fr *find.FindResult) bool {

		if name == "" {
			return true
		}

		if fr != nil && strings.Contains(fr.Name, name) || strings.Contains(fr.Extra, name) {
			return true
		}
		/*
			if fr != nil && fr.Source == "reddit" {
				return true
			}
		*/
		fmt.Printf("RESULT for : [%s] is not correct: %+v\n", name, fr)
		return false
	}

	// set the filter func to find the specific files
	filterOpt = filter.NewGenericOptions(root).SetFilterFunc(fFunc)

	// this will not compile if the package did not implement the inteface.
	switch findType {
	case "social":
		repository = social.NewRepository()
	default:
		repository = filesystem.NewRepository()
	}

	fmt.Printf("Looking for Name: (%s) in Respository: %s\n", name, findType)
	// find the result in the repository
	fAbj = find.New(repository)
	fr, err := fAbj.Do(filterOpt)

	if err != nil {
		panic(err)
	}

	//spew.Dump(fr)
	//panic("!!!")
	if err := ui.Display(fr); err != nil {
		panic(err)
	}

}

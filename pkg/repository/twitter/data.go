// log into twitter, gather all likes,
// store likes in a datastore,
// search the datastore for link content when implementing the grep
// note the difference between find and search
package twitter

import (
	"github.com/dathan/go-find-hexagonal/pkg/ctwitter"
	"github.com/dathan/go-find-hexagonal/pkg/find"
	"github.com/davecheney/errors"
)

// Respository is the local struct that implements the find interface
type Repository struct {
}

// Respository  returns the struct
func NewRepository() *Repository {
	ret := &Repository{}
	return ret
}

// Implements the repository interface
func (f *Repository) Find(fo find.FilterOptions) (find.FindResults, error) {

	allTweets, err := ctwitter.NewService().Favorites.List(fo.GetStart())
	if err != nil {
		return nil, errors.Annotate(err, "Repository.Find()")
	}

	// return the filtered results.
	var res find.FindResults
	for _, tweet := range allTweets {

		path := "n/a"
		lenUrls := len(tweet.Entities.Urls)

		if lenUrls > 0 {
			path = ""
			for pos, url := range tweet.Entities.Urls {
				sep := ""
				if pos > 0 && lenUrls > 1 && pos != lenUrls-1 {
					sep = ":"
				}
				path = path + sep + url.ExpandedURL
			}
		}

		// Find Keywords in the Result - this is equivalent to finding a pattern in the filename
		fResult := find.FindResult{Name: tweet.FullText, Path: path}
		if fo.GetFilterFunc()(&fResult) {
			res = append(res, fResult)
		}
	}

	return res, nil
}

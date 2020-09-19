// log into twitter, gather all likes,
// store likes in a datastore,
// search the datastore for link content when implementing the grep
// note the difference between find and search
package social

import (
	"fmt"
	"sync"

	"github.com/dathan/go-find-hexagonal/pkg/find"
	reddit "github.com/dathan/go-find-hexagonal/pkg/social/reddit"
	twitter "github.com/dathan/go-find-hexagonal/pkg/social/twitter"
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

// Implements the repository interface that concurrently pulls from Reddit and Twitter
func (f *Repository) Find(fo find.FilterOptions) (find.FindResults, error) {

	var wg sync.WaitGroup
	wg.Add(1)

	results := make(chan find.FindResults, 2)
	fatalErrors := make(chan error)
	wgDone := make(chan bool)
	// wait in a routine and close the results so this method returns
	go func() {
		wg.Wait()
		close(results)
		close(wgDone)
	}()
	/*
		go func(results chan find.FindResults) {
			defer wg.Done()
			res, err := f.twitterFind(fo)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				fatalErrors <- err
			}

			results <- res
		}(results)
	*/
	go func(results chan find.FindResults) {

		defer wg.Done()
		res, err := f.redditFind(fo)

		if err != nil {
			fatalErrors <- err
			return
		}

		fmt.Printf("about to write res to results\n")
		results <- res
		fmt.Printf("all done\n")

	}(results)

	select {
	case err := <-fatalErrors:
		close(fatalErrors)
		return nil, err
	case <-wgDone:
		fmt.Printf("all wgAreDone")
		var allResults find.FindResults
		for res := range results {
			if res != nil {
				allResults = append(allResults, res...)
			}
		}
		fmt.Printf("Returning results: %d\n", len(allResults))
		return allResults, nil
	}

	panic(errors.New("impossible error"))
}

func (f *Repository) redditFind(fo find.FilterOptions) (find.FindResults, error) {

	service, err := reddit.NewService()
	if err != nil {
		return nil, errors.Wrap(err, errors.Errorf("redditFind service failed () - %s ", err))
	}

	// only works because my handle is the same between platforms
	upVotes, err := service.Favorites.CacheList(fo.GetStart())

	if err != nil {
		return nil, err
	}

	filteredResults := find.FindResults{}

	for _, vote := range upVotes {

		fResult := find.FindResult{
			Name:   vote.Title,
			Path:   vote.FullPermalink(),
			Extra:  vote.Selftext,
			Source: "reddit",
		}

		// filterAllResults
		if fo.GetFilterFunc()(&fResult) {
			filteredResults = append(filteredResults, fResult)
		}

	}

	fmt.Printf("filteredResults: %d\n", len(filteredResults))
	return filteredResults, nil
}
func (f *Repository) twitterFind(fo find.FilterOptions) (find.FindResults, error) {
	service, err := twitter.NewService()

	if err != nil {
		return nil, err
	}

	allTweets, err := service.Favorites.CacheList(fo.GetStart())
	if err != nil {
		return nil, errors.Annotate(err, "Repository.Find()")
	}

	// return the filtered results.
	var res find.FindResults
	for _, tweet := range allTweets {

		path := "https://twitter.com/i/web/status/" + tweet.IDStr
		lenUrls := len(tweet.Entities.Urls)

		if lenUrls > 0 {
			for _, urlStruct := range tweet.Entities.Urls {
				if path != urlStruct.ExpandedURL {
					path += "\n\n" + urlStruct.ExpandedURL
				}
			}
		}

		t, err := tweet.CreatedAtTime()
		if err != nil {
			return nil, err
		}

		// Find Keywords in the Result - this is equivalent to finding a pattern in the filename
		fResult := find.FindResult{
			Name:      tweet.Text,
			Extra:     tweet.FullText,
			Path:      path,
			CreatedAt: t.Unix(),
		}

		if fo.GetFilterFunc()(&fResult) {
			res = append(res, fResult)
		}
	}

	return res, nil
}

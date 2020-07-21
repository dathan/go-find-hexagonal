// log into twitter, gather all likes,
// store likes in a datastore,
// search the datastore for link content when implementing the grep
// note the difference between find and search
package twitter

import (
	"fmt"
	"os"

	"github.com/dathan/go-find-hexagonal/pkg/find"
	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
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
	// TODO: pull in data from twitter, index it, find the link in another thread
	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("TWITTER_CONSUMER_KEY"),
		ClientSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	twitterHandle := fo.GetStart()

	// Twitter client
	client := twitter.NewClient(httpClient)
	var allTweets []twitter.Tweet
	var maxID int64 = int64(0)

	// keep paginating until there are no more results from maxID
	for {

		options := &twitter.FavoriteListParams{ScreenName: twitterHandle, IncludeEntities: twitter.Bool(true)}
		// set the pagination
		if maxID != 0 {
			options.MaxID = maxID
		}

		// have tweets
		tweets, _, err := client.Favorites.List(options)
		if tweets == nil || err != nil {
			if err != nil {
				fmt.Printf("What is this error: %s\n", err)
			}
			break
		}

		//return all the tweets to the result filter.
		allTweets = append(allTweets, tweets...)

		// if there is a cycle
		if maxID > 0 && maxID == tweets[len(tweets)-1].ID {
			break
		}

		maxID = tweets[len(tweets)-1].ID // returns the result less than this id so we don't need to subtract
	}

	// return the filtered results.
	var res find.FindResults
	for _, tweet := range allTweets {

		path := "n/a"
		lenUrls := len(tweet.Entities.Urls)

		if lenUrls > 0 {
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

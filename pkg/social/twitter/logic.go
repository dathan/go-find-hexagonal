package twitter

import (
	"context"
	"os"

	"github.com/dathan/go-find-hexagonal/pkg/cache"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/clientcredentials"
)

const MAX_PAGES = 20

// NewService returns the store and error
func NewService() (*Store, error) {

	t := &Store{Favorites: &Bookmarks{}}

	key := os.Getenv("TWITTER_CONSUMER_KEY")
	secret := os.Getenv("TWITTER_CONSUMER_SECRET")
	if len(key) == 0 || len(secret) == 0 {
		return nil, errors.New("twitter.NewService() - invalid environment")
	}

	config := &clientcredentials.Config{
		ClientID:     key,
		ClientSecret: secret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	// http.Client will automatically authorize Requests
	httpClient := config.Client(context.Background())
	// Twitter client
	client := twitter.NewClient(httpClient)
	t.client = client
	t.Favorites.parent = t

	return t, nil
}

func (tFav *Bookmarks) setupCache() error {

	if tFav.cache == nil {
		favCache, err := cache.NewService()
		if err != nil {
			return err
		}
		tFav.cache = favCache

	}
	return nil

}

func (tFav *Bookmarks) CacheList(twitterHandle string) (Tweets, error) {

	if err := tFav.setupCache(); err != nil {
		return nil, err
	}

	allTweets := Tweets{}
	key := "Favorites.twitter." + twitterHandle
	ok, err := tFav.cache.Get(key, &allTweets)
	if err == nil && ok {
		//fmt.Printf("RETURNING CACHE: %v\n", allTweets)
		return allTweets, nil
	}

	allTweets, err = tFav.List(twitterHandle)

	if err != nil {
		return nil, err
	}

	if err := tFav.cache.Set(key, allTweets); err != nil {
		return nil, err
	}

	return allTweets, nil

}

// List
func (tFav *Bookmarks) List(twitterHandle string) (Tweets, error) {

	maxID := int64(0)
	page := 0
	allTweets := Tweets{}

	for {
		page++
		options := &twitter.FavoriteListParams{ScreenName: twitterHandle, SinceID: 0, Count: 200, IncludeEntities: twitter.Bool(true)}
		// set the pagination
		if maxID != 0 {
			options.MaxID = maxID
		}

		// have tweets
		tweets, _, err := tFav.parent.client.Favorites.List(options)
		if tweets == nil || err != nil {
			if err != nil {
				if len(allTweets) == 0 {
					return nil, err
				}
			}
			break
		}

		if len(tweets) == 0 {
			break
		}

		tSize := len(tweets)
		//fmt.Printf("ArrayRange [%d]: %d N: %d LEN: %d\n", maxID, tweets[0].ID, tweets[tSize-1].ID, len(tweets))

		//return all the tweets to the result filter.
		allTweets = append(allTweets, tweets...)

		// if there is a cycle
		if maxID > 0 && maxID == tweets[tSize-1].ID {
			break
		}

		maxID = tweets[tSize-1].ID

		if page > MAX_PAGES {
			break
		}
	}

	return allTweets, nil
}

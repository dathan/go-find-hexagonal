package ctwitter

import (
	"context"
	"fmt"
	"os"

	"github.com/dathan/go-find-hexagonal/pkg/cache"
	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2/clientcredentials"
)

type TwitterFavorites struct {
	parent *Twitter
	cache  cache.Store
}

type Twitter struct {
	Favorites *TwitterFavorites
	client    *twitter.Client
}

// type declaration to abstract the type
type Tweets []twitter.Tweet

type Favorites interface {
	All(string) (Tweets, error)
}

func NewService() *Twitter {
	t := &Twitter{Favorites: &TwitterFavorites{}}
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("TWITTER_CONSUMER_KEY"),
		ClientSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	// http.Client will automatically authorize Requests
	httpClient := config.Client(context.Background())
	// Twitter client
	client := twitter.NewClient(httpClient)
	t.client = client
	t.Favorites.parent = t

	return t
}

func (tFav *TwitterFavorites) setupCache() error {

	if tFav.cache == nil {
		favCache, err := cache.NewService()
		if err != nil {
			return err
		}
		tFav.cache = favCache

	}
	return nil

}

func (tFav *TwitterFavorites) List(twitterHandle string) (Tweets, error) {

	key := "Favorites." + twitterHandle

	if err := tFav.setupCache(); err != nil {
		return nil, err
	}

	allTweets := Tweets{}
	ok, err := tFav.cache.Get(key, allTweets)
	if err == nil && ok {
		fmt.Printf("RETURNING CACHE")
		return allTweets, nil
	}

	maxID := int64(0)
	page := 0
	for {
		page++
		options := &twitter.FavoriteListParams{ScreenName: twitterHandle, IncludeEntities: twitter.Bool(true)}
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

		//return all the tweets to the result filter.
		allTweets = append(allTweets, tweets...)

		// if there is a cycle
		if maxID > 0 && maxID == tweets[len(tweets)-1].ID {
			break
		}

		maxID = tweets[len(tweets)-1].ID // returns the result less than this id so we don't need to subtract
		fmt.Printf("PAGE: %d Size: %d \n", page, len(allTweets))

		if page > 2 {
			break
		}
	}

	if err := tFav.cache.Set(key, allTweets); err != nil {
		return nil, err
	}

	return allTweets, nil
}

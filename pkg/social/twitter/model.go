package twitter

import (
	"github.com/dathan/go-find-hexagonal/pkg/cache"
	"github.com/dghubble/go-twitter/twitter"
)

// Everything that was liked, favorited conceptually to me when you mark a post that's a bookmark or a type of endourcement for memory.
type Bookmarks struct {
	parent *Store // Store contain's the client to perform the actions (Twitter for instance)
	cache  cache.Store
}

type Store struct {
	Favorites *Bookmarks
	client    *twitter.Client
}

// type declaration to abstract the type
type Tweets []twitter.Tweet

// Favorites interface
type Favorites interface {
	List(string) (Tweets, error)
	CacheList(string) (Tweets, error)
}

// let's define the Service Interface to return the store or an error
type Service interface {
	NewService(*Store, error)
}

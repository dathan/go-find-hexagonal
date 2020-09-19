package reddit

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"

	"github.com/dathan/geddit"
	"github.com/dathan/go-find-hexagonal/pkg/cache"
)

const MAX_PAGES = 20

type UpVotes struct {
	parent *Client
	cache  cache.Store
}

type Client struct {
	Favorites *UpVotes
	client    *geddit.OAuthSession
}

type Favorites interface {
	List(string) (Upvoted, error)
	CacheList(string) (Upvoted, error)
}

type Service interface {
	NewService() (*Client, error)
}

// Upvoted type definition
type Upvoted []*geddit.Submission

func NewService() (*Client, error) {
	t := &Client{Favorites: &UpVotes{}}

	key := os.Getenv("REDDIT_CONSUMER_KEY")
	secret := os.Getenv("REDDIT_CONSUMER_SECRET")
	botUser := os.Getenv("REDDIT_USERNAME")
	botPass := os.Getenv("REDDIT_PASSWORD")

	checkInputs := []string{key, secret, botUser, botPass}

	for _, in := range checkInputs {
		if len(in) == 0 {
			return nil, errors.New("reddit.NewService() - invalid environment - " + fmt.Sprintf("%+v", checkInputs))
		}
	}

	o, err := geddit.NewOAuthSession(
		key,
		secret,
		"UpVoteFinder",
		"http://localhost/oauth",
	)

	if err != nil {
		return nil, err
	}

	// Create new auth token for confidential clients (personal scripts/apps).
	err = o.LoginAuth(botUser, botPass)
	if err != nil {
		return nil, errors.Wrap(err, "reddit.NewService() - invalid environment")
	}

	t.client = o
	t.Favorites.parent = t

	// Ready to make API calls!
	return t, nil
}

func (tFav *UpVotes) setupCache() error {

	if tFav.cache == nil {
		favCache, err := cache.NewService()
		if err != nil {
			return err
		}
		tFav.cache = favCache

	}
	return nil

}

func (tFav *UpVotes) CacheList(handle string) (Upvoted, error) {

	if err := tFav.setupCache(); err != nil {
		return nil, err
	}

	cacheKey := "reddit.Upvotes." + handle
	var upvoted Upvoted = make(Upvoted, 0)

	if ok, err := tFav.cache.Get(cacheKey, &upvoted); err == nil && ok {
		return upvoted, nil
	}

	upvoted, err := tFav.List(handle)
	if err != nil {
		return nil, err
	}

	if err := tFav.cache.Set(cacheKey, upvoted); err != nil {
		return nil, err
	}

	return upvoted, nil

}

func (tFav *UpVotes) List(handle string) (Upvoted, error) {

	var sort geddit.PopularitySort = "top"
	var params geddit.ListingOptions = geddit.ListingOptions{Show: "all", Count: 200, Limit: 200}
	var allUpvoted Upvoted
	page := 0
	for {
		page++
		if page > 1 {
			fmt.Printf("Page: %d\n", page)
		}

		subs, err := tFav.parent.client.Upvoted(handle, sort, params)
		if err != nil {
			return allUpvoted, errors.Wrap(err, "reddit.UpVotes.List() - ")
		}

		if len(subs) == 0 {
			fmt.Println("No more subs")
			break
		}

		if len(allUpvoted) > 0 && subs[0].FullID == allUpvoted[len(subs)-1].FullID {
			fmt.Printf("Pulling the same information from last: %v\n", subs)
			break
		}

		before := subs[len(subs)-1].FullID
		if len(params.After) > 0 && before == params.After {
			fmt.Printf("BREAK ON BEFORE: %s\n", before)
			break
		}

		params.After = before
		allUpvoted = append(allUpvoted, subs...)

	}

	spew.Dump(allUpvoted)

	return allUpvoted, nil
}

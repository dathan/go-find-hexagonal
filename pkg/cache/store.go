package cache

import (
	"time"

	"github.com/davecheney/errors"
	"github.com/philippgille/gokv/bbolt"
)

// we could use a type alias instead of redefine the interface but we are not migrating a package
// type aliases e.g. type s = r are good to inherit the type so you can factor it out of your code for instance
// type declaraation e.g. type F string use when you want to add additional behavor to the concrete type
type Store interface {
	Set(k string, v interface{}) error
	Get(k string, v interface{}) (found bool, err error)
	Delete(k string) error
	Close() error
}

const cacheDelta = 3600

type StoreCache struct {
	db bbolt.Store
}

type TimeDeltaCache struct {
	orig  interface{}
	start time.Time
}

func NewService() (Store, error) {
	// use bolt db as the db
	options := bbolt.Options{}
	bdb, err := bbolt.NewStore(options) // use default settings todo fix
	if err != nil {
		return nil, errors.Annotate(err, "cache.NewService()")
	}

	lCache := &StoreCache{
		db: bdb,
	}

	defer lCache.Close()

	return lCache, nil

}

func (cache *StoreCache) Set(k string, v interface{}) error {
	lru := TimeDeltaCache{
		orig:  v,
		start: time.Now(),
	}
	return cache.db.Set(k, lru)
}

func (cache *StoreCache) Close() error {
	return cache.db.Close()
}

func (cache *StoreCache) Delete(k string) error {
	return cache.db.Delete(k)
}

// note v is an pointer so the type is changed to the value of the object past assuming its the same type
func (cache *StoreCache) Get(k string, v interface{}) (bool, error) {
	// TODO put LRU logic here for the get (we will purge on selects :) )
	lru := TimeDeltaCache{}

	ok, err := cache.db.Get(k, lru)
	if err != nil {
		return false, err
	}

	if time.Now().Unix()-lru.start.Unix() > cacheDelta {
		return false, nil
	}

	return ok, err

}

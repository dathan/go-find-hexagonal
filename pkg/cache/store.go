package cache

import (
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

type LRUCache struct {
	db bbolt.Store
}

func NewService() (Store, error) {
	// use bolt db as the db
	options := bbolt.Options{}
	bdb, err := bbolt.NewStore(options) // use default settings todo fix
	if err != nil {
		return nil, errors.Annotate(err, "cache.NewService()")
	}

	lCache := &LRUCache{
		db: bdb,
	}

	defer lCache.Close()

	return lCache, nil

}

func (cache *LRUCache) Set(k string, v interface{}) error {
	return cache.db.Set(k, v)
}

func (cache *LRUCache) Close() error {
	return cache.db.Close()
}

func (cache *LRUCache) Delete(k string) error {
	return cache.db.Delete(k)
}

// note v is an pointer so the type is changed to the value of the object past assuming its the same type
func (cache *LRUCache) Get(k string, v interface{}) (bool, error) {
	// TODO put LRU logic here for the get (we will purge on selects :) )
	return cache.db.Get(k, v)
}

package cache

import (
	"reflect"
	"sync"
	"time"

	"github.com/davecheney/errors"
	"github.com/philippgille/gokv/badgerdb"
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
	db  *badgerdb.Store
	mux sync.Mutex
}

var singleBadgerDB *badgerdb.Store

type TimeDeltaCache struct {
	Orig  interface{}
	Start time.Time
}

func NewService() (Store, error) {
	// use bolt db as the db
	var mux sync.Mutex
	mux.Lock()
	defer mux.Unlock()

	options := badgerdb.Options{}
	if singleBadgerDB == nil {
		bdb, err := badgerdb.NewStore(options) // use default settings todo fix
		if err != nil {
			return nil, errors.Annotate(err, "cache.NewService()")
		}
		singleBadgerDB = &bdb
	}

	lCache := &StoreCache{
		db: singleBadgerDB,
	}

	//defer lCache.Close()

	return lCache, nil

}

func (cache *StoreCache) Set(k string, v interface{}) error {

	cache.mux.Lock()
	defer cache.mux.Unlock()

	tC := &TimeDeltaCache{v, time.Now()}

	return cache.db.Set(k, tC)
}

func (cache *StoreCache) Close() error {
	return cache.db.Close()
}

func (cache *StoreCache) Delete(k string) error {
	cache.mux.Lock()
	defer cache.mux.Unlock()
	return cache.db.Delete(k)
}

// note v is an pointer so the type is changed to the value of the object past assuming its the same type
func (cache *StoreCache) Get(k string, v interface{}) (bool, error) {
	cache.mux.Lock()
	defer cache.mux.Unlock()
	// TODO put LRU logic here for the get (we will purge on selects :) )
	tDCache := &TimeDeltaCache{
		Orig: v, // ensure the v is here so its the right type that gets set by reference
	}

	ok, err := cache.db.Get(k, tDCache)
	if err != nil {
		return false, err
	}

	if time.Now().Unix()-tDCache.Start.Unix() > cacheDelta {
		clear(&v)
		return false, nil
	}

	return ok, err

}

func clear(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

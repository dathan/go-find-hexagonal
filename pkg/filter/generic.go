// Package filter implements the required interface for find
package filter

import (
	"github.com/dathan/go-find-hexagonal/pkg/find"
)

// GenericOptions is a generic option list.
type GenericOptions struct {
	config     map[string]string
	path       string
	filterFunc find.FilterFunc
}

// NewGenericOptions return the generic options
func NewGenericOptions(start string) *GenericOptions {
	return &GenericOptions{
		config: make(map[string]string),
		path:   start,
	}
}

// GetStart return the starting point
func (fo *GenericOptions) GetStart() string {
	return fo.path
}

// Set the options map untile we come up with a better method
func (fo *GenericOptions) Set(k, v string) {
	fo.config[k] = v
}

// SetFilterFunc set the filter we are going to use if any
func (fo *GenericOptions) SetFilterFunc(ff find.FilterFunc) find.FilterOptions {
	fo.filterFunc = ff
	return fo
}

// GetFilterFunc return the set filterfunc or nil
func (fo *GenericOptions) GetFilterFunc() find.FilterFunc {
	if fo.filterFunc != nil {
		return fo.filterFunc
	}

	return find.FilterFunc(func(fr *find.FindResult) bool {
		return fr != nil
	})
}

// Get the options from the map
func (fo *GenericOptions) Get(k string) string {
	ret := fo.config[k]
	return ret
}

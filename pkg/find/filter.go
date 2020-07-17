package find

type FilterFunc func(*FindResult) bool

// define the filter interface to only have results which are filtered
type FilterOptions interface {
	GetStart() string
	SetFilterFunc(FilterFunc) FilterOptions
	GetFilterFunc() FilterFunc
}

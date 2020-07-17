package find

// The interace to how storage lookups are done
type Repository interface {
	Find(fo FilterOptions) (FindResults, error)
}

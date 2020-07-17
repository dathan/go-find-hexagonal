package find

type Find interface {
	Do(fo FilterOptions) (FindResults, error)
}

type find struct {
	repository Repository
}

func New(res Repository) Find {
	ret := &find{repository: res}
	return ret
}

// use the repository to find
func (f *find) Do(fo FilterOptions) (FindResults, error) {
	return f.repository.Find(fo)
}

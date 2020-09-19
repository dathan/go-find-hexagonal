package find

// FindResult the generic results from a repository
type FindResult struct {
	Name      string
	CreatedAt int64
	Path      string
	Extra     string
	Source    string
}

// type declaration to short hand the result
type FindResults []FindResult

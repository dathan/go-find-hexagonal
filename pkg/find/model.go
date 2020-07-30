package find

type FindResult struct {
	Name      string
	CreatedAt int64
	Path      string
	Extra     string
	// TODO Add more things here
}

type FindResults []FindResult

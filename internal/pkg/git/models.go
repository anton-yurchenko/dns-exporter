package vcs

// Project represents a git repository
type Project struct {
	Name        string
	AuthorName  string
	AuthorEmail string
	Remote      *Origin
}

// Origin represents a remote git origin
type Origin struct {
	URL    string
	Branch string
}

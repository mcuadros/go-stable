package gopkg

import "gopkg.in/sourcegraph/go-vcsurl.v1"

type Repository struct {
	URL string
}

func NewRepository(url string) (*Repository, error) {
	info, err := vcsurl.Parse(url)
	if err != nil {
		return nil, err
	}

	return &Repository{URL: info.CloneURL}, nil
}

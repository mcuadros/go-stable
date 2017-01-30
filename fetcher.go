package stable

import (
	"io"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/protocol/packp"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type Fetcher struct {
	pkg     *Package
	service transport.UploadPackSession
	auth    transport.AuthMethod
}

func NewFetcher(p *Package, auth transport.AuthMethod) *Fetcher {

	s, _ := http.DefaultClient.NewUploadPackSession(p.Repository, auth)

	return &Fetcher{pkg: p, service: s}
}

func (f *Fetcher) Versions() (Versions, error) {
	info, err := f.service.AdvertisedReferences()
	if err != nil {
		return nil, err
	}

	refs, err := info.AllReferences()
	if err != nil {
		return nil, err
	}

	return NewVersions(refs), nil
}

func (f *Fetcher) Fetch(w io.Writer, ref *plumbing.Reference) (written int64, err error) {
	req := packp.NewUploadPackRequest()
	req.Wants = []plumbing.Hash{ref.Hash()}

	r, err := f.service.UploadPack(req)
	if err != nil {
		return 0, err
	}

	return io.Copy(w, r)
}

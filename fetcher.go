package gopkg

import (
	"io"

	"gopkg.in/src-d/go-git.v4/clients/common"
	"gopkg.in/src-d/go-git.v4/clients/http"
	"gopkg.in/src-d/go-git.v4/core"
)

type Fetcher struct {
	service common.GitUploadPackService
	pkg     *Package
	auth    common.AuthMethod
}

func NewFetcher(p *Package, auth common.AuthMethod) *Fetcher {
	s := http.NewGitUploadPackService(p.Repository)
	s.SetAuth(auth)

	return &Fetcher{pkg: p, service: s}
}

func (f *Fetcher) Versions() (Versions, error) {
	if err := f.service.Connect(); err != nil {
		return nil, err
	}

	info, err := f.service.Info()
	if err != nil {
		return nil, err
	}

	return NewVersions(info.Refs), nil
}

func (f *Fetcher) Fetch(w io.Writer, ref *core.Reference) (written int64, err error) {
	if err := f.service.Connect(); err != nil {
		return 0, err
	}

	req := &common.GitUploadPackRequest{}
	req.Want(ref.Hash())

	r, err := f.service.Fetch(req)
	if err != nil {
		return 0, err
	}

	return io.Copy(w, r)
}

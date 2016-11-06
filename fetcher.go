package gopkg

import (
	"io"

	"github.com/mxk/go-flowrate/flowrate"
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

	return NewVersions(info), nil
}

func (f *Fetcher) Fetch(w io.Writer, ref *core.Reference) (*flowrate.Status, error) {
	if err := f.service.Connect(); err != nil {
		return nil, err
	}

	req := &common.GitUploadPackRequest{}
	req.Want(ref.Hash())

	r, err := f.service.Fetch(req)
	if err != nil {
		return nil, err
	}

	flow := flowrate.NewReader(r, -1)
	defer flow.Close()

	_, err = io.Copy(w, flow)
	status := flow.Status()
	return &status, err
}

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
	service := http.NewGitUploadPackService(p.Repository)
	service.SetAuth(auth)

	return &Fetcher{pkg: p, auth: auth, service: service}
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
	//req.Depth = 1

	r, err := f.service.Fetch(req)
	if err != nil {
		return nil, err
	}

	flow := flowrate.NewReader(r, -1)
	_, err = io.Copy(w, flow)
	flow.Close()

	status := flow.Status()
	return &status, err
}

package gopkg

import (
	"io"

	"github.com/mxk/go-flowrate/flowrate"
	"gopkg.in/src-d/go-git.v2"
	"gopkg.in/src-d/go-git.v2/clients/common"
	"gopkg.in/src-d/go-git.v2/core"
)

const defaultBranch = "refs/heads/master"

type Fetcher struct {
	pkg    *Package
	remote *git.Remote
	auth   common.AuthMethod
}

func NewFetcher(p *Package, auth common.AuthMethod) *Fetcher {
	return &Fetcher{pkg: p, auth: auth}
}

func (f *Fetcher) Info() (*common.GitUploadPackInfo, error) {
	var err error
	f.remote, err = git.NewAuthenticatedRemote(f.pkg.Repository.CloneURL, f.auth)
	if err != nil {
		return nil, err
	}

	if err := f.remote.Connect(); err != nil {
		return nil, err
	}

	info := f.remote.Info()
	v := NewVersions(info).Match(f.pkg.Repository.Rev)
	if v == nil {
		return nil, ErrVersionNotFound
	}

	info.Head = v.Hash
	info.Refs = map[string]core.Hash{defaultBranch: v.Hash}

	return info, nil
}

func (f *Fetcher) Fetch(w io.Writer) (*flowrate.Status, error) {
	if _, err := f.Info(); err != nil {
		return nil, err
	}

	r, err := f.remote.FetchDefaultBranch()
	if err != nil {
		return nil, err
	}

	flow := flowrate.NewReader(r, -1)
	_, err = io.Copy(w, flow)
	flow.Close()

	status := flow.Status()
	return &status, err
}

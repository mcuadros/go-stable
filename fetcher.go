package gopkg

import (
	"io"

	"github.com/mxk/go-flowrate/flowrate"
	"gopkg.in/src-d/go-git.v2"
	"gopkg.in/src-d/go-git.v2/clients/common"
	"gopkg.in/src-d/go-git.v2/core"
)

type Fetcher struct {
	isConnected bool
	pkg         *Package
	remote      *git.Remote
	auth        common.AuthMethod
}

func NewFetcher(p *Package, auth common.AuthMethod) *Fetcher {
	return &Fetcher{pkg: p, auth: auth}
}

func (f *Fetcher) Versions() (Versions, error) {
	if err := f.connect(); err != nil {
		return nil, err
	}

	info := f.remote.Info()
	return NewVersions(info), nil
}

func (f *Fetcher) Fetch(w io.Writer, ref core.Hash) (*flowrate.Status, error) {
	if err := f.connect(); err != nil {
		return nil, err
	}

	req := &common.GitUploadPackRequest{}
	req.Want(ref)

	r, err := f.remote.Fetch(req)
	if err != nil {
		return nil, err
	}

	flow := flowrate.NewReader(r, -1)
	_, err = io.Copy(w, flow)
	flow.Close()

	status := flow.Status()
	return &status, err
}

func (f *Fetcher) connect() error {
	if f.isConnected {
		return nil
	}

	defer func() { f.isConnected = true }()

	var err error
	f.remote, err = git.NewAuthenticatedRemote(f.pkg.Repository.CloneURL, f.auth)
	if err != nil {
		return err
	}

	if err := f.remote.Connect(); err != nil {
		return err
	}

	return nil
}

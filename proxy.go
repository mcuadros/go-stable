package gopkg

import (
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/src-d/go-git.v2/core"
)

const DefaultBranch = "refs/heads/master"

type Proxy struct {
	Notifiers struct {
		InfoRefs      func(*http.Request, error)
		GitUploadPack func(*http.Request, error)
		RawSaved      func(*http.Request, error)
	}
}

func NewProxy() *Proxy {
	return &Proxy{}
}

func (p *Proxy) Start(addr string) error {
	http.HandleFunc("/", p.handler)
	return http.ListenAndServe(addr, nil)
}

func (p *Proxy) StartTLS(addr, certFile, keyFile string) error {
	http.HandleFunc("/", p.handler)
	return http.ListenAndServeTLS(addr, certFile, keyFile, nil)
}

func (p *Proxy) handler(w http.ResponseWriter, r *http.Request) {
	var err error
	switch {
	case strings.HasSuffix(r.URL.Path, "/info/refs"):
		err = p.doUploadPackInfoResponse(w, r)
	case strings.HasSuffix(r.URL.Path, "/git-upload-pack"):
		err = p.doUploadPackResponse(w, r)
	default:
		err = p.defaultHandler(w, r)
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

var template = `<html><head><meta name="go-import" content="gop.kg/%s git http://gop.kg/%#[1]s"></head><body></body></html>`

func (p *Proxy) defaultHandler(w http.ResponseWriter, r *http.Request) error {
	if r.FormValue("go-get") != "1" {
		return fmt.Errorf("invalid request: %s", r.URL.Path)
	}

	w.Header().Set("Content-Type", "text/html")
	_, err := fmt.Fprintf(w, template, r.URL.Path[1:])
	return err
}

func (p *Proxy) doUploadPackInfoResponse(w http.ResponseWriter, r *http.Request) error {
	url := strings.Replace(r.URL.Path, "/info/refs", "", 1)

	repository, err := NewRepository("https://github.com" + url)
	if err != nil {
		return err
	}

	fetcher := NewFetcher(repository)
	info, err := fetcher.Info()
	if err != nil {
		return err
	}

	fmt.Println("--->", info.String())

	info.Head = "refs/heads/master"
	info.Refs = map[string]core.Hash{
		"refs/heads/master": info.Refs["refs/heads/master"],
	}

	w.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")
	_, err = w.Write(info.Bytes())
	return err
}

func (p *Proxy) doUploadPackResponse(w http.ResponseWriter, r *http.Request) error {
	url := strings.Replace(r.URL.Path, "/git-upload-pack", "", 1)

	repository, err := NewRepository("https://github.com" + url)
	if err != nil {
		return err
	}

	fetcher := NewFetcher(repository)

	w.Write([]byte("0008NAK\n"))
	if _, err := fetcher.Fetch(w); err != nil {
		return err
	}

	return nil
}

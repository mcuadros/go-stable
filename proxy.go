package gopkg

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gopkg.in/src-d/go-git.v2/clients/common"
	githttp "gopkg.in/src-d/go-git.v2/clients/http"
	"gopkg.in/src-d/go-git.v2/core"
)

const DefaultBranch = "refs/heads/master"

var errUncatchedRequest = errors.New("uncatched request")

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

type Context struct {
	Package *Package
	*gin.Context
}

func (p *Proxy) Handle(c *gin.Context) error {
	pkg, err := NewPackageFromRequest(c.Request)
	if err != nil {
		return err
	}

	c.Set("package", pkg)
	return p.do(&Context{pkg, c})
}

func (p *Proxy) do(c *Context) error {
	var err error
	switch {
	case strings.HasSuffix(c.Request.URL.Path, "/info/refs"):
		err = p.doUploadPackInfoResponse(c)
	case strings.HasSuffix(c.Request.URL.Path, "/git-upload-pack"):
		err = p.doUploadPackResponse(c)
	default:
		err = p.defaultHandler(c)
	}

	return p.handleError(c, err)
}

var template = `<html><head><meta name="go-import" content="%s git http://%#[1]s"></head><body></body></html>`

func (p *Proxy) defaultHandler(c *Context) error {
	if c.Query("go-get") != "1" {
		return errUncatchedRequest
	}

	c.Header("Content-Type", "text/html")
	_, err := fmt.Fprintf(c.Writer, template, c.Package.Name)

	return err
}

func (p *Proxy) doUploadPackInfoResponse(c *Context) error {
	fetcher := NewFetcher(c.Package, p.getAuth(c))
	info, err := fetcher.Info()
	if err != nil {
		return err
	}

	info.Head = info.Refs["refs/heads/master"]
	info.Refs = map[string]core.Hash{
		"refs/heads/master": info.Refs["refs/heads/master"],
	}

	c.Header("Content-Type", "application/x-git-upload-pack-advertisement")
	c.String(200, info.String())

	return nil
}

func (p *Proxy) doUploadPackResponse(c *Context) error {
	c.Header("Content-Type", "application/x-git-upload-pack-result")
	if _, err := c.Writer.WriteString("0008NAK\n"); err != nil {
		return err
	}

	fetcher := NewFetcher(c.Package, p.getAuth(c))
	if _, err := fetcher.Fetch(c.Writer); err != nil {
		return err
	}

	return nil
}

func (p *Proxy) handleError(c *Context, err error) error {
	if err == errUncatchedRequest {
		return nil
	}

	c.Abort()
	if err, ok := err.(*core.PermanentError); ok {
		if err.Err == common.NotFoundErr {
			return p.handleNotFoundError(c, err)
		}
	}

	return err
}

func (p *Proxy) handleNotFoundError(c *Context, err error) error {
	if !p.isAuth(c) {
		p.requireAuth(c)
	} else {
		c.AbortWithError(404, err)
	}

	return nil
}

func (p *Proxy) requireAuth(c *Context) bool {
	if _, _, ok := c.Request.BasicAuth(); ok {
		return true
	}

	c.Header("WWW-Authenticate", `Basic realm="GoPkg"`)
	c.String(401, "401 Unauthorized")
	c.Abort()

	return false
}

func (p *Proxy) getAuth(c *Context) *githttp.BasicAuth {
	username, password, _ := c.Request.BasicAuth()

	return githttp.NewBasicAuth(username, password)
}

func (p *Proxy) isAuth(c *Context) bool {
	_, _, ok := c.Request.BasicAuth()

	return ok
}

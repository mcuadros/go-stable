package gopkg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/src-d/go-git.v2/clients/common"
	"gopkg.in/src-d/go-git.v2/clients/http"
	"gopkg.in/src-d/go-git.v2/core"
)

const defaultBranch = "refs/heads/master"

var (
	ErrUncatchedRequest = errors.New("uncatched request")
	ErrVersionNotFound  = errors.New("version not found")
)

type Proxy struct {
	Notifiers struct {
		InfoRefs      func(*Context, error)
		GitUploadPack func(*Context, error)
		RawSaved      func(*Context, error)
	}
}

type Context struct {
	Package *Package
	*gin.Context
}

func NewProxy() *Proxy {
	return &Proxy{}
}

func (p *Proxy) Handle(c *gin.Context) error {
	ctx := &Context{nil, c}
	pkg, err := NewPackageFromRequest(c.Request)
	if err != nil {
		return p.handleError(ctx, err)
	}

	c.Set("package", pkg)
	ctx.Package = pkg
	return p.do(ctx)
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

var gogetTemplate = `<html><head><meta name="go-import" content="%s git http://%#[1]s"></head><body></body></html>`

func (p *Proxy) defaultHandler(c *Context) error {
	if c.Query("go-get") != "1" {
		return ErrUncatchedRequest
	}

	c.Header("Content-Type", "text/html")
	_, err := fmt.Fprintf(c.Writer, gogetTemplate, c.Package.Name)

	return err
}

func (p *Proxy) doUploadPackInfoResponse(c *Context) error {
	fetcher := NewFetcher(c.Package, p.getAuth(c))
	v, err := p.getVersion(fetcher, c.Package)
	if err != nil {
		return err
	}

	info := common.NewGitUploadPackInfo()
	info.Head = v.Hash
	info.Refs = map[string]core.Hash{defaultBranch: v.Hash}
	info.Capabilities.Set("symref", "HEAD:"+defaultBranch)

	c.Header("Content-Type", "application/x-git-upload-pack-advertisement")
	c.String(200, info.String())

	return nil
}

func (p *Proxy) doUploadPackResponse(c *Context) error {
	fetcher := NewFetcher(c.Package, p.getAuth(c))
	v, err := p.getVersion(fetcher, c.Package)
	if err != nil {
		return err
	}

	c.Header("Content-Type", "application/x-git-upload-pack-result")
	if _, err := c.Writer.WriteString("0008NAK\n"); err != nil {
		return err
	}

	if _, err := fetcher.Fetch(c.Writer, v.Hash); err != nil {
		return err
	}

	return nil
}

func (p *Proxy) getVersion(f *Fetcher, pkg *Package) (*Version, error) {
	versions, err := f.Versions()
	if err != nil {
		return nil, err
	}

	v := versions.Match(pkg.Repository.Rev)
	if v == nil {
		return nil, ErrVersionNotFound
	}

	return v, nil
}

func (p *Proxy) handleError(c *Context, err error) error {
	if err == ErrUncatchedRequest || err == ErrInvalidRequest {
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

func (p *Proxy) getAuth(c *Context) *http.BasicAuth {
	username, password, _ := c.Request.BasicAuth()

	return http.NewBasicAuth(username, password)
}

func (p *Proxy) isAuth(c *Context) bool {
	_, _, ok := c.Request.BasicAuth()

	return ok
}

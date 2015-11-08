package gopkg

import (
	"net/http"
	"net/url"
	"testing"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git.v2/clients/common"
	"gopkg.in/src-d/go-git.v2/core"
)

func Test(t *testing.T) { TestingT(t) }

type SuiteCommon struct{}

var _ = Suite(&SuiteCommon{})

func (s *SuiteCommon) TestNewRepositoryFromRequestSubdomain(c *C) {
	p := s.buildPackageFromRequest(c, "http://foo.gop.kg/qux@baz", nil)
	c.Assert(p.Name, Equals, "foo.gop.kg/qux@baz")
	c.Assert(p.Repository.CloneURL, Equals, "git://github.com/foo/qux.git")
	c.Assert(p.Repository.Rev, Equals, "baz")

	p = s.buildPackageFromRequest(c, "http://foo.gop.kg/qux@baz/info/refs?service=git-upload-pack", nil)
	c.Assert(p.Name, Equals, "foo.gop.kg/qux@baz")
	c.Assert(p.Repository.CloneURL, Equals, "git://github.com/foo/qux.git")
	c.Assert(p.Repository.Rev, Equals, "baz")
}

func (s *SuiteCommon) TestNewRepositoryFromRequestPath(c *C) {
	UrlMode = Path
	defer func() { UrlMode = Subdomain }()

	p := s.buildPackageFromRequest(c, "http://gop.kg/foo/qux@master", nil)
	c.Assert(p.Name, Equals, "gop.kg/foo/qux@master")
	c.Assert(p.Repository.CloneURL, Equals, "git://github.com/foo/qux.git")
	c.Assert(p.Repository.Rev, Equals, "master")

	p = s.buildPackageFromRequest(c, "http://gop.kg/foo/qux@master", nil)
	c.Assert(p.Name, Equals, "gop.kg/foo/qux@master")
	c.Assert(p.Repository.CloneURL, Equals, "git://github.com/foo/qux.git")
	c.Assert(p.Repository.Rev, Equals, "master")

	p = s.buildPackageFromRequest(c, "http://gop.kg/foo/qux@master/bar", nil)
	c.Assert(p.Name, Equals, "gop.kg/foo/qux@master/bar")
	c.Assert(p.Repository.CloneURL, Equals, "git://github.com/foo/qux.git")
	c.Assert(p.Repository.Rev, Equals, "master")

	s.buildPackageFromRequest(c, "http://gop.kg/foo/qux/bar@master", ErrInvalidRequest)
}

func (s *SuiteCommon) buildPackageFromRequest(c *C, reqURL string, expectedErr error) *Package {
	url, _ := url.Parse(reqURL)
	req := &http.Request{}
	req.URL = url

	p, err := NewPackageFromRequest(req)
	c.Assert(err, Equals, expectedErr)

	return p
}

func (s *SuiteCommon) TestNewVersionHead(c *C) {
	v := NewVersion("refs/heads/qux", core.NewHash(""))
	c.Assert(v.Name, Equals, "qux")
	c.Assert(v.Type, Equals, Head)
}

func (s *SuiteCommon) TestNewVersionTag(c *C) {
	v := NewVersion("refs/tags/foo", core.NewHash(""))
	c.Assert(v.Name, Equals, "foo")
	c.Assert(v.Type, Equals, Tag)
}

func (s *SuiteCommon) TestNewVersionAnnotatedTag(c *C) {
	v := NewVersion("refs/tags/foo^{}", core.NewHash(""))
	c.Assert(v.Name, Equals, "foo")
	c.Assert(v.Type, Equals, AnnotatedTag)
}

func (s *SuiteCommon) TestNewVersions(c *C) {
	info := &common.GitUploadPackInfo{}
	info.Refs = map[string]core.Hash{
		"refs/heads/master":  core.NewHash(""),
		"refs/tags/v1.0.0":   core.NewHash(""),
		"refs/tags/1.1.2":    core.NewHash(""),
		"refs/tags/1.1.2^{}": core.NewHash(""),
		"refs/tags/1.1.3":    core.NewHash(""),
		"refs/tags/1.1.3^{}": core.NewHash(""),
		"refs/tags/v1.0.3":   core.NewHash(""),
		"refs/tags/v2.0.3":   core.NewHash(""),
	}

	v := NewVersions(info)
	c.Assert(v.Match("1.1").Ref, Equals, "refs/tags/1.1.3^{}")
	c.Assert(v.Match("1.1.2").Ref, Equals, "refs/tags/1.1.2^{}")
	c.Assert(v.Match("2").Ref, Equals, "refs/tags/v2.0.3")
	c.Assert(v.Match("master").Ref, Equals, "refs/heads/master")
	c.Assert(v.Match("foo"), IsNil)
}

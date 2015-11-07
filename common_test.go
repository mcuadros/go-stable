package gopkg

import (
	"net/http"
	"net/url"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type SuiteCommon struct{}

var _ = Suite(&SuiteCommon{})

func (s *SuiteCommon) TestNewRepositoryFromRequest(c *C) {
	p := s.buildPackageFromRequest(c, "http://foo.gop.kg/qux")
	c.Assert(p.Name, Equals, "foo.gop.kg/qux")
	c.Assert(p.Repository.CloneURL, Equals, "git://github.com/foo/qux.git")

	p = s.buildPackageFromRequest(c, "http://gop.kg/foo/qux")
	c.Assert(p.Name, Equals, "gop.kg/foo/qux")
	c.Assert(p.Repository.CloneURL, Equals, "git://github.com/foo/qux.git")

	p = s.buildPackageFromRequest(c, "http://gop.kg/foo/qux/bar")
	c.Assert(p.Name, Equals, "gop.kg/foo/qux/bar")
	c.Assert(p.Repository.CloneURL, Equals, "git://github.com/foo/qux.git")

	p = s.buildPackageFromRequest(c, "http://foo.gop.kg/qux/info/refs?service=git-upload-pack")
	c.Assert(p.Name, Equals, "foo.gop.kg/qux")
	c.Assert(p.Repository.CloneURL, Equals, "git://github.com/foo/qux.git")
}

func (s *SuiteCommon) buildPackageFromRequest(c *C, reqURL string) *Package {
	url, _ := url.Parse(reqURL)
	req := &http.Request{}
	req.URL = url

	p, err := NewPackageFromRequest(req)
	c.Assert(err, IsNil)

	return p
}

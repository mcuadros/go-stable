package stable

import (
	"bytes"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/client/common"
)

type FetcherSuite struct{}

var _ = Suite(&FetcherSuite{})

func (s *FetcherSuite) TestVersions(c *C) {
	pkg := &Package{}
	pkg.Repository, _ = common.NewEndpoint("https://github.com/git-fixtures/basic")

	f := NewFetcher(pkg, nil)
	versions, err := f.Versions()
	c.Assert(err, IsNil)
	c.Assert(versions, HasLen, 2)
}

func (s *FetcherSuite) TestFetch(c *C) {
	pkg := &Package{}
	pkg.Repository, _ = common.NewEndpoint("https://github.com/git-fixtures/basic")

	f := NewFetcher(pkg, nil)

	ref := plumbing.NewReferenceFromStrings("foo", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5")

	buf := bytes.NewBuffer(nil)
	n, err := f.Fetch(buf, ref)
	c.Assert(err, IsNil)
	c.Assert(n, Equals, int64(85374))
	c.Assert(buf.Len(), Equals, 85374)
}

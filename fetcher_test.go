package gopkg

import (
	"bytes"

	. "gopkg.in/check.v1"
)

type FetcherSuite struct{}

var _ = Suite(&FetcherSuite{})

func (s *FetcherSuite) TestInfo(c *C) {
	pkg := &Package{}
	pkg.Repository.CloneURL = "https://github.com/tyba/git-fixture"
	pkg.Repository.Rev = "master"

	f := NewFetcher(pkg, nil)
	info, err := f.Info()

	c.Assert(err, IsNil)
	c.Assert(info.Head.String(), Equals, "6ecf0ef2c2dffb796033e5a02219af86ec6584e5")
}

func (s *FetcherSuite) TestFetch(c *C) {
	pkg := &Package{}
	pkg.Repository.CloneURL = "https://github.com/tyba/git-fixture"
	pkg.Repository.Rev = "master"

	f := NewFetcher(pkg, nil)

	buf := bytes.NewBuffer(nil)
	status, err := f.Fetch(buf)
	c.Assert(err, IsNil)
	c.Assert(status.Active, Equals, false)
	c.Assert(status.Bytes, Equals, int64(85374))
	c.Assert(buf.Len(), Equals, 85374)
}

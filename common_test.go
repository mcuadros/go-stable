package stable

import (
	"testing"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func Test(t *testing.T) { TestingT(t) }

type SuiteCommon struct{}

var _ = Suite(&SuiteCommon{})

func (s *SuiteCommon) TestNewVersions(c *C) {
	refs := make(memory.ReferenceStorage, 0)
	refs.SetReference(plumbing.NewHashReference("refs/heads/master", plumbing.NewHash("")))
	refs.SetReference(plumbing.NewHashReference("refs/tags/v1.0.0", plumbing.NewHash("")))
	refs.SetReference(plumbing.NewHashReference("refs/tags/1.1.2", plumbing.NewHash("")))
	refs.SetReference(plumbing.NewHashReference("refs/tags/1.1.3", plumbing.NewHash("")))
	refs.SetReference(plumbing.NewHashReference("refs/tags/v1.0.3", plumbing.NewHash("")))
	refs.SetReference(plumbing.NewHashReference("refs/tags/v2.0.3", plumbing.NewHash("")))
	refs.SetReference(plumbing.NewHashReference("refs/tags/v4.0.0-rc1", plumbing.NewHash("")))

	v := NewVersions(refs)
	c.Assert(v.BestMatch("v0").Name().String(), Equals, "refs/heads/master")
	c.Assert(v.BestMatch("v1.1").Name().String(), Equals, "refs/tags/1.1.3")
	c.Assert(v.BestMatch("1.1").Name().String(), Equals, "refs/tags/1.1.3")
	c.Assert(v.BestMatch("1.1.2").Name().String(), Equals, "refs/tags/1.1.2")
	c.Assert(v.BestMatch("2").Name().String(), Equals, "refs/tags/v2.0.3")
	c.Assert(v.BestMatch("4").Name().String(), Equals, "refs/tags/v4.0.0-rc1")
	c.Assert(v.BestMatch("master").Name().String(), Equals, "refs/heads/master")
	c.Assert(v.BestMatch("foo"), IsNil)

	refs.SetReference(plumbing.NewHashReference("refs/tags/v0.0.0", plumbing.NewHash("")))

	v = NewVersions(refs)
	c.Assert(v.BestMatch("v0").Name().String(), Equals, "refs/tags/v0.0.0")
}

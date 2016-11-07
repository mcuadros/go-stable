package stable

import (
	"testing"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git.v4/core"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func Test(t *testing.T) { TestingT(t) }

type SuiteCommon struct{}

var _ = Suite(&SuiteCommon{})

func (s *SuiteCommon) TestNewVersions(c *C) {
	refs := make(memory.ReferenceStorage, 0)
	refs.SetReference(core.NewHashReference("refs/heads/master", core.NewHash("")))
	refs.SetReference(core.NewHashReference("refs/tags/v1.0.0", core.NewHash("")))
	refs.SetReference(core.NewHashReference("refs/tags/1.1.2", core.NewHash("")))
	refs.SetReference(core.NewHashReference("refs/tags/1.1.3", core.NewHash("")))
	refs.SetReference(core.NewHashReference("refs/tags/v1.0.3", core.NewHash("")))
	refs.SetReference(core.NewHashReference("refs/tags/v2.0.3", core.NewHash("")))
	refs.SetReference(core.NewHashReference("refs/tags/v4.0.0-rc1", core.NewHash("")))

	v := NewVersions(refs)
	c.Assert(v.BestMatch("v1.1").Name().String(), Equals, "refs/tags/1.1.3")
	c.Assert(v.BestMatch("1.1").Name().String(), Equals, "refs/tags/1.1.3")
	c.Assert(v.BestMatch("1.1.2").Name().String(), Equals, "refs/tags/1.1.2")
	c.Assert(v.BestMatch("2").Name().String(), Equals, "refs/tags/v2.0.3")
	c.Assert(v.BestMatch("4").Name().String(), Equals, "refs/tags/v4.0.0-rc1")
	c.Assert(v.BestMatch("master").Name().String(), Equals, "refs/heads/master")
	c.Assert(v.BestMatch("foo"), IsNil)
}

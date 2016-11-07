package stable

import (
	"fmt"

	"github.com/mcuadros/go-version"
	"gopkg.in/src-d/go-git.v4/clients/common"
	"gopkg.in/src-d/go-git.v4/core"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// Package represent a golang package
type Package struct {
	Name       string
	Repository common.Endpoint
	Constrain  string
	Versions   Versions
}

type Versions map[string]*core.Reference

func NewVersions(refs memory.ReferenceStorage) Versions {
	versions := make(Versions, 0)
	for _, ref := range refs {
		if !ref.IsTag() && !ref.IsBranch() {
			continue
		}

		versions[ref.Name().Short()] = ref
	}

	return versions
}

func (v Versions) Match(needed string) []*core.Reference {
	c := newConstrain(needed)

	var names []string
	for _, ref := range v {
		name := ref.Name().Short()
		if c.Match(version.Normalize(name)) {
			names = append(names, name)
		}
	}

	version.Sort(names)
	var matched []*core.Reference
	for n := len(names) - 1; n >= 0; n-- {
		matched = append(matched, v[names[n]])
	}

	return matched
}

func (v Versions) BestMatch(needed string) *core.Reference {
	if version, ok := v[needed]; ok {
		return version
	}

	matched := v.Match(needed)
	if len(matched) == 0 {
		return nil
	}

	return matched[0]
}

func (v Versions) Mayor() map[string]*core.Reference {
	output := make(map[string]*core.Reference, 0)
	for i := 0; i < 100; i++ {
		mayor := fmt.Sprintf("v%d", i)
		if m := v.BestMatch(mayor); m != nil {
			output[mayor] = m
		}
	}

	return output
}

func newConstrain(needed string) *version.ConstraintGroup {
	if needed[0] == 'v' && needed[1] >= 28 && needed[1] <= 57 {
		needed = needed[1:]
	}

	return version.NewConstrainGroupFromString(needed + ".*")
}

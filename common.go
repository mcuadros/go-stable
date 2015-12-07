package gopkg

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mcuadros/go-version"
	"gopkg.in/sourcegraph/go-vcsurl.v1"
	"gopkg.in/src-d/go-git.v2/clients/common"
	"gopkg.in/src-d/go-git.v2/core"
)

type urlMode int

const (
	Subdomain urlMode = 1
	Path      urlMode = 2
)

var (
	VersionSeparator  = "!"
	UrlMode           = Path
	ErrInvalidRequest = errors.New("invalid request")

	ignoredPrefixes = []string{"/git-upload-pack", "/info/refs"}
)

type PackageName string

func (p PackageName) Base() string {
	parts := splitByVersionSeparator(p)
	return parts[0]
}

func (p PackageName) Version() string {
	parts := splitByVersionSeparator(p)
	if len(parts) < 2 {
		return ""
	}

	return parts[1]
}

func (p PackageName) Change(v *Version) PackageName {
	return PackageName(p.Base() + VersionSeparator + v.Name)
}

func splitByVersionSeparator(n PackageName) []string {
	parts := strings.Split(string(n), VersionSeparator)
	return parts
}

// Package represent a golang package
type Package struct {
	Name       PackageName
	Repository vcsurl.RepoInfo
	Versions   Versions
}

// NewPackage returns a new instance of Package
func NewPackage(p PackageName, url string) (*Package, error) {
	info, err := vcsurl.Parse(url)
	if err != nil {
		return nil, err
	}

	return &Package{Name: p, Repository: *info}, nil
}

// NewPackageFromRequest returns a new instance of Package using the data from
// a http request
func NewPackageFromRequest(r *http.Request) (*Package, error) {
	username := getUsername(r)
	repository, revision := getRepository(r)

	if username == "" || repository == "" || revision == "" {
		return nil, ErrInvalidRequest
	}

	return NewPackage(
		getPackageName(r),
		fmt.Sprintf("https://github.com/%s/%s#%s", username, repository, revision),
	)
}

func getUsername(r *http.Request) string {
	switch UrlMode {
	case Subdomain:
		names := strings.Split(getHost(r), ".")
		if len(names) > 2 {
			return names[0]
		}
	case Path:
		path := strings.Split(getPath(r), "/")
		if len(path) > 1 {
			return path[0]
		}
	}

	return ""
}

func getRepository(r *http.Request) (repository, version string) {
	switch UrlMode {
	case Subdomain:
		path := getPath(r)
		if strings.Index(path, "/") != -1 {
			return
		}

		return splitRepository(path)
	case Path:
		path := strings.Split(getPath(r), "/")
		if len(path) == 1 {
			return
		}

		return splitRepository(path[1])
	}

	return
}

func splitRepository(s string) (repository, version string) {
	parts := strings.Split(s, VersionSeparator)
	if len(parts) == 1 {
		return s, ""
	}

	return parts[0], parts[1]
}

func getPackageName(r *http.Request) PackageName {
	return PackageName(fmt.Sprintf("%s/%s", getHost(r), getPath(r)))
}

func getHost(r *http.Request) string {
	if r.Host != "" {
		return r.Host
	}

	return r.URL.Host
}

func getPath(r *http.Request) string {
	path := r.URL.Path[1:]
	if strings.Count(path, VersionSeparator) > 1 {
		return ""
	}

	for _, prefix := range ignoredPrefixes {
		path = strings.Replace(path, prefix, "", -1)
	}

	return path
}

type VersionType int

const (
	Head         VersionType = 1
	Tag          VersionType = 2
	AnnotatedTag VersionType = 3
)

type Version struct {
	Name string
	Ref  string
	Hash core.Hash
	Type VersionType
}

func NewVersion(ref string, h core.Hash) *Version {
	v := &Version{Ref: ref, Hash: h}

	switch {
	case strings.HasPrefix(ref, "refs/tags"):
		v.Type = Tag
	case strings.HasPrefix(ref, "refs/heads"):
		v.Type = Head
	default:
		return nil
	}

	v.Name = ref[strings.LastIndex(ref, "/")+1:]
	if strings.HasSuffix(v.Name, "^{}") {
		v.Type = AnnotatedTag
		v.Name = v.Name[:len(v.Name)-3]
	}

	return v
}

type Versions map[string]*Version

func NewVersions(info *common.GitUploadPackInfo) Versions {
	versions := make(Versions, 0)
	aVersions := make(Versions, 0)
	for ref, hash := range info.Refs {
		v := NewVersion(ref, hash)
		if v == nil {
			continue
		}

		if v.Type == AnnotatedTag {
			aVersions[v.Name] = v
		} else {
			versions[v.Name] = v
		}
	}

	for _, v := range aVersions {
		versions[v.Name] = v
	}

	return versions
}

func (v Versions) Match(needed string) []*Version {
	c := newConstrain(needed)

	var names []string
	for _, ver := range v {
		if c.Match(version.Normalize(ver.Name)) {
			names = append(names, ver.Name)
		}
	}

	version.Sort(names)
	var matched []*Version
	for n := len(names) - 1; n >= 0; n-- {
		matched = append(matched, v[names[n]])
	}

	return matched
}

func (v Versions) BestMatch(needed string) *Version {
	if version, ok := v[needed]; ok {
		return version
	}

	matched := v.Match(needed)
	if len(matched) == 0 {
		return nil
	}

	return matched[0]
}

func (v Versions) Mayor() map[string]*Version {
	output := make(map[string]*Version, 0)
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

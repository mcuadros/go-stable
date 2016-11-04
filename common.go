package gopkg

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mcuadros/go-version"
	"gopkg.in/src-d/go-git.v4/clients/common"
	"gopkg.in/src-d/go-git.v4/core"
)

type urlMode int

const (
	Subdomain urlMode = 1
	Path      urlMode = 2
)

var (
	VersionSeparator  = "@"
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

	parts = strings.Split(parts[1], "/")
	return parts[0]
}

func (p PackageName) Root() string {
	return p.Base() + VersionSeparator + p.Version()
}

func (p PackageName) Change(v *core.Reference) PackageName {
	return PackageName(p.Base() + VersionSeparator + v.Name().Short())
}

func splitByVersionSeparator(n PackageName) []string {
	parts := strings.Split(string(n), VersionSeparator)
	return parts
}

// Package represent a golang package
type Package struct {
	Name       PackageName
	Repository common.Endpoint
	Condition  string
	Versions   Versions
}

// NewPackage returns a new instance of Package
func NewPackage(p PackageName, url string) (*Package, error) {

	return nil, nil
}

// NewPackageFromRequest returns a new instance of Package using the data from
// a http request
func NewPackageFromRequest(r *http.Request) (*Package, error) {
	username := getUsername(r)
	repository, revision := getRepository(r)
	name := getPackageName(r)

	url, err := getURL(username, repository)
	if err != nil {
		return nil, err
	}

	if username == "" || repository == "" || revision == "" {
		return nil, ErrInvalidRequest
	}

	return &Package{
		Name:       name,
		Repository: url,
		Condition:  revision,
	}, nil
}

func getURL(username, repository string) (endpoint common.Endpoint, err error) {
	return common.NewEndpoint(fmt.Sprintf(
		"https://github.com/%s/%s", username, repository,
	))
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

type Versions map[string]*core.Reference

func NewVersions(info *common.GitUploadPackInfo) Versions {
	versions := make(Versions, 0)
	for _, ref := range info.Refs {
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
		fmt.Println("name, name", name)
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

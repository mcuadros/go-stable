package gopkg

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/sourcegraph/go-vcsurl.v1"
)

var (
	ErrInvalidRequest = errors.New("invalid request")

	ignoredPrefixes = []string{"/git-upload-pack", "/info/refs"}
)

// Package represent a golang package
type Package struct {
	Name       string
	Repository vcsurl.RepoInfo
}

// NewPackage returns a new instance of Package
func NewPackage(pkgName, url string) (*Package, error) {
	info, err := vcsurl.Parse(url)
	if err != nil {
		return nil, err
	}

	return &Package{Name: pkgName, Repository: *info}, nil
}

// NewPackageFromRequest returns a new instance of Package using the data from
// a http request
func NewPackageFromRequest(r *http.Request) (*Package, error) {
	username := getUsername(r)
	repository := getRepository(r)

	if username == "" || repository == "" {
		return nil, ErrInvalidRequest
	}

	return NewPackage(
		getPackageName(r),
		fmt.Sprintf("https://github.com/%s/%s", username, repository),
	)
}

func getUsername(r *http.Request) string {
	names := strings.Split(getHost(r), ".")
	if len(names) > 2 {
		return names[0]
	}

	path := strings.Split(getPath(r), "/")
	if len(path) > 1 {
		return path[0]
	}

	return ""
}

func getRepository(r *http.Request) string {
	path := strings.Split(getPath(r), "/")
	names := strings.Split(getHost(r), ".")
	if len(names) > 2 {
		return path[0]
	}

	if len(path) < 2 {
		return ""
	}

	return path[1]
}

func getPackageName(r *http.Request) string {
	return fmt.Sprintf("%s/%s", getHost(r), getPath(r))
}

func getHost(r *http.Request) string {
	if r.Host != "" {
		return r.Host
	}

	return r.URL.Host
}

func getPath(r *http.Request) string {
	path := r.URL.Path[1:]
	for _, prefix := range ignoredPrefixes {
		path = strings.Replace(path, prefix, "", -1)
	}

	return path
}

package gopkg

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"gopkg.in/src-d/go-git.v4/clients/common"
	githttp "gopkg.in/src-d/go-git.v4/clients/http"
	"gopkg.in/src-d/go-git.v4/core"
	"gopkg.in/src-d/go-git.v4/formats/packp/pktline"

	"github.com/gorilla/mux"
)

type Server struct {
	Base string
	Host string

	Default struct {
		Server       string
		Organazation string
		Repository   string
		Constrain    string
	}

	r *mux.Router
}

func NewServer() *Server {
	s := &Server{}

	s.Base = "/{org:[a-z0-9-]+}/{repository:[a-z0-9-]+}@{version:v[0-9.]+}"
	s.Host = "example.com"
	s.Default.Server = "github.com"

	return s
}

func (s *Server) Run(addr ...string) error {
	s.buildRouter()
	http.Handle("/", s.r)

	log.Println("Listening...")
	http.ListenAndServe(addr[0], nil)
	return nil
}

func (s *Server) RunTLS(addr, certFile, keyFile string) error {
	s.buildRouter()
	http.Handle("/", s.r)

	log.Println("Listening...")
	http.ListenAndServeTLS(addr, certFile, keyFile, nil)
	return nil
}

func (s *Server) buildRouter() {
	s.r = mux.NewRouter()

	s.r.HandleFunc(s.Base, s.doMetaImportResponse).Methods("GET").Name("base")
	s.r.HandleFunc(s.Base, s.doMetaImportResponse).Methods("GET").Queries("go-get", "1")
	s.r.HandleFunc(path.Join(s.Base, "/info/refs"), s.doUploadPackInfoResponse).Methods("GET")
	s.r.HandleFunc(path.Join(s.Base, "/git-upload-pack"), s.doUploadPackResponse)
}

const (
	ServerKey       = "server"
	OrganizationKey = "org"
	RepositoryKey   = "repository"
	ConstraintKey   = "version"
)

func (s *Server) doMetaImportResponse(w http.ResponseWriter, r *http.Request) {
	pkg := s.buildPackage(r)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, gogetTemplate, pkg.Name)
}

func (s *Server) doUploadPackInfoResponse(w http.ResponseWriter, r *http.Request) {
	pkg := s.buildPackage(r)

	fetcher := NewFetcher(pkg, getAuth(r))
	ref, err := s.getVersion(fetcher, pkg)
	if err != nil {
		panic(err)
	}

	ref = s.mutateTagToBranch(ref)
	info := s.buildGitUploadPackInfo(ref)

	w.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")
	w.Write(info.Bytes())
}

func (s *Server) getVersion(f *Fetcher, pkg *Package) (*core.Reference, error) {
	versions, err := f.Versions()
	if err != nil {
		return nil, err
	}

	v := versions.BestMatch(pkg.Constrain)
	if v == nil {
		return nil, ErrVersionNotFound
	}

	return v, nil
}

// we mutate the tag into a branch to avoid detached branches and allow to the
// user now the current tag selected
func (s *Server) mutateTagToBranch(ref *core.Reference) *core.Reference {
	if ref.IsBranch() {
		return ref
	}

	branch := core.ReferenceName("refs/heads/" + ref.Name().Short())
	return core.NewHashReference(branch, ref.Hash())

}

func (s *Server) buildGitUploadPackInfo(ref *core.Reference) *common.GitUploadPackInfo {
	info := common.NewGitUploadPackInfo()
	info.Refs.SetReference(ref)
	info.Refs.SetReference(core.NewSymbolicReference(core.HEAD, ref.Name()))
	info.Capabilities.Set("symref", "HEAD:"+ref.Name().String())

	return info
}

func (s *Server) buildPackage(r *http.Request) *Package {
	params := mux.Vars(r)
	server := getOrDefault(params, ServerKey, s.Default.Server)
	organization := getOrDefault(params, OrganizationKey, s.Default.Organazation)
	repository := getOrDefault(params, RepositoryKey, s.Default.Repository)

	name, _ := s.r.Get("base").URL(
		"server", server,
		"org", organization,
		"repository", repository,
		"version", params[ConstraintKey],
	)

	return &Package{
		Name:       path.Join(s.Host, name.String()),
		Repository: s.buildEndpoint(server, organization, repository),
		Constrain:  params[ConstraintKey],
	}
}

func (s *Server) buildEndpoint(server, orgnization, repository string) common.Endpoint {
	e, err := common.NewEndpoint(fmt.Sprintf(
		"https://%s/%s/%s", server, orgnization, repository,
	))

	if err != nil {
		panic(fmt.Sprintf("unreachable: %s", err.Error()))
	}

	return e
}

func (s *Server) doUploadPackResponse(w http.ResponseWriter, r *http.Request) {
	pkg := s.buildPackage(r)
	fetcher := NewFetcher(pkg, getAuth(r))
	ref, err := s.getVersion(fetcher, pkg)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/x-git-upload-pack-result")

	pkt := pktline.NewEncoder(w)
	if err := pkt.EncodeString("NAK\n"); err != nil {
		panic(err)
	}

	_, err = fetcher.Fetch(w, ref)
	if err != nil {
		panic(err)
	}
}

func getOrDefault(m map[string]string, key, def string) string {
	if v, ok := m[key]; ok {
		return v
	}

	return def
}

func getAuth(r *http.Request) *githttp.BasicAuth {
	username, password, _ := r.BasicAuth()

	return githttp.NewBasicAuth(username, password)
}

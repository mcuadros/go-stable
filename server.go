package stable

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

const DefaultBaseRoute = "/{org:[a-z0-9-]+}/{repository:[a-z0-9-/]+}.{version:v[0-9.]+}"

type Server struct {
	http.Server
	r *mux.Router

	BaseRoute string
	Host      string
	Default   struct {
		Server       string
		Organization string
		Repository   string
	}
}

func NewDefaultServer(host string) *Server {
	s := NewServer(DefaultBaseRoute, host)
	s.Default.Server = "github.com"

	return s
}

func NewServer(base, host string) *Server {
	s := &Server{
		BaseRoute: base,
		Host:      host,
	}

	s.buildRouter()
	return s
}

func (s *Server) ListenAndServe() error {
	panic("ListenAndServer, is not supported try ListenAndServeTLS")
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	return s.Server.ListenAndServeTLS(certFile, keyFile)
}

func (s *Server) buildRouter() {
	s.r = mux.NewRouter()
	s.r.HandleFunc("/", s.doRootRedirect).Methods("GET").Name("base")
	s.r.HandleFunc(path.Join(s.BaseRoute, "/info/refs"), s.doUploadPackInfoResponse).Methods("GET")
	s.r.HandleFunc(path.Join(s.BaseRoute, "/git-upload-pack"), s.doUploadPackResponse).Methods("POST")
	s.r.HandleFunc(path.Join(s.BaseRoute, "/{subpkg:.+}"), s.doMetaImportResponse).Methods("GET").Queries("go-get", "1")
	s.r.HandleFunc(path.Join(s.BaseRoute, "/{subpkg:.+}"), s.doPackageRedirect).Methods("GET")
	s.r.HandleFunc(s.BaseRoute, s.doMetaImportResponse).Methods("GET").Queries("go-get", "1")
	s.r.HandleFunc(s.BaseRoute, s.doPackageRedirect).Methods("GET").Name("base")
	s.Handler = s.r
}

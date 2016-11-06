package gopkg

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

const DefaultBaseRoute = "/{org:[a-z0-9-]+}/{repository:[a-z0-9-/]+}@{version:v[0-9.]+}"

type Server struct {
	http.Server
	r *mux.Router

	Base    string
	Host    string
	Default struct {
		Server       string
		Organazation string
		Repository   string
		Constrain    string
	}
}

func NewDefaultServer(host string) *Server {
	s := NewServer(DefaultBaseRoute, host)
	s.Default.Server = "github.com"

	return s
}

func NewServer(base, host string) *Server {
	return &Server{
		Base: base,
		Host: host,
	}
}

func (s *Server) ListenAndServe() error {
	panic("ListenAndServer, is not supported try ListenAndServeTLS")
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	s.buildRouter()
	return s.Server.ListenAndServeTLS(certFile, keyFile)
}

func (s *Server) buildRouter() {
	s.r = mux.NewRouter()
	s.r.HandleFunc(s.Base, s.doMetaImportResponse).Methods("GET").Name("base")
	s.r.HandleFunc(s.Base, s.doMetaImportResponse).Methods("GET").Queries("go-get", "1")
	s.r.HandleFunc(path.Join(s.Base, "/info/refs"), s.doUploadPackInfoResponse).Methods("GET")
	s.r.HandleFunc(path.Join(s.Base, "/git-upload-pack"), s.doUploadPackResponse).Methods("POST")

	s.Handler = s.r
}

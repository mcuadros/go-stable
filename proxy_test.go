package stable

import . "gopkg.in/check.v1"
import "net/http"
import "net/http/httptest"

import "io/ioutil"

type ProxySuite struct{}

var _ = Suite(&ProxySuite{})

func (s *ProxySuite) TestDoMetaImportResponse(c *C) {
	r, _ := http.NewRequest("GET", "http://foo.bar/git-fixtures/releases.v1?go-get=1", nil)
	s.doTestDoMetaImportResponse(c, r)
}

func (s *ProxySuite) TestDoMetaImportResponseSubpackage(c *C) {
	r, _ := http.NewRequest("GET", "http://foo.bar/git-fixtures/releases/subpackage.v1?go-get=1", nil)
	s.doTestDoMetaImportResponse(c, r)
}

func (s *ProxySuite) doTestDoMetaImportResponse(c *C, r *http.Request) {
	w := httptest.NewRecorder()

	server := NewDefaultServer("foo.bar")
	server.buildRouter()
	server.Handler.ServeHTTP(w, r)

	response := w.Result()
	body, err := ioutil.ReadAll(response.Body)
	c.Assert(err, IsNil)
	c.Assert(response.StatusCode, Equals, 200)

	c.Assert(string(body), Equals, ""+
		"<html>\n"+
		"\t\t<head>\n"+
		"\t\t\t<meta name=\"go-import\" content=\"foo.bar/git-fixtures/releases.v1 git https://foo.bar/git-fixtures/releases.v1\">\n"+
		"\t\t</head>\n"+
		"\t\t<body></body>\n"+
		"\t</html>",
	)

	c.Assert(response.Header.Get("Content-Type"), Equals, "text/html")
}

func (s *ProxySuite) TestDoUploadPackInfoResponse(c *C) {
	r, _ := http.NewRequest("GET", "http://foo.bar/git-fixtures/releases.v1/info/refs", nil)
	w := httptest.NewRecorder()

	server := NewDefaultServer("foo.bar")
	server.buildRouter()
	server.Handler.ServeHTTP(w, r)

	response := w.Result()
	body, err := ioutil.ReadAll(response.Body)
	c.Assert(err, IsNil)
	c.Assert(string(body), Equals, ""+
		"001e# service=git-upload-pack\n"+
		"0000005096f2c336f6aec28963719fb42513b88dfd709d09 HEAD\x00symref=HEAD:refs/heads/v1.0.0\n"+
		"003f96f2c336f6aec28963719fb42513b88dfd709d09 refs/heads/v1.0.0\n"+
		"0000",
	)

	c.Assert(response.StatusCode, Equals, 200)
	c.Assert(response.Header.Get("Content-Type"), Equals, "application/x-git-upload-pack-advertisement")
}

func (s *ProxySuite) TestDoUploadPackInfoResponsePrivate(c *C) {
	r, _ := http.NewRequest("GET", "http://foo.bar/git-fixtures/private.v1/info/refs", nil)
	w := httptest.NewRecorder()

	server := NewDefaultServer("foo.bar")
	server.buildRouter()
	server.Handler.ServeHTTP(w, r)

	response := w.Result()
	c.Assert(response.StatusCode, Equals, 401)
}

func (s *ProxySuite) TestDoUploadPackResponse(c *C) {
	r, _ := http.NewRequest("POST", "http://foo.bar/git-fixtures/releases.v1/git-upload-pack", nil)
	w := httptest.NewRecorder()

	server := NewDefaultServer("foo.bar")
	server.buildRouter()
	server.Handler.ServeHTTP(w, r)

	response := w.Result()
	body, err := ioutil.ReadAll(response.Body)
	c.Assert(err, IsNil)
	c.Assert(len(body), Equals, 1152)

	c.Assert(response.StatusCode, Equals, http.StatusOK)
	c.Assert(response.Header.Get("Content-Type"), Equals, "application/x-git-upload-pack-result")
}

func (s *ProxySuite) TestDoRootRedirect(c *C) {
	r, _ := http.NewRequest("GET", "http://foo.bar/", nil)
	w := httptest.NewRecorder()

	server := NewDefaultServer("foo.bar")
	server.Default.Server = "qux.baz"
	server.Default.Organization = "foo"

	server.buildRouter()
	server.Handler.ServeHTTP(w, r)

	response := w.Result()
	body, err := ioutil.ReadAll(response.Body)
	c.Assert(err, IsNil)
	c.Assert(len(body), Equals, 42)

	c.Assert(response.StatusCode, Equals, http.StatusFound)
	c.Assert(response.Header.Get("Location"), Equals, "https://qux.baz/foo")
}

func (s *ProxySuite) TestDoPackageRedirect(c *C) {
	r, _ := http.NewRequest("GET", "http://foo.bar/org/repository.v1", nil)
	w := httptest.NewRecorder()

	server := NewDefaultServer("foo.bar")
	server.buildRouter()
	server.Handler.ServeHTTP(w, r)

	response := w.Result()
	body, err := ioutil.ReadAll(response.Body)
	c.Assert(err, IsNil)
	c.Assert(len(body), Equals, 56)

	c.Assert(response.StatusCode, Equals, http.StatusFound)
	c.Assert(response.Header.Get("Location"), Equals, "https://github.com/org/repository")
}

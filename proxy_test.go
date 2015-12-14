package gopkg

import (
	"bufio"
	"bytes"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"
)

type ProxySuite struct{}

var _ = Suite(&ProxySuite{})

func (s *ProxySuite) TestDefaultHandlerGoGet(c *C) {
	ctx := s.newMockContext("GET", "https://foo.gop.kg/bar!baz?go-get=1")

	p := &Proxy{}

	err := p.Handle(ctx)
	c.Assert(err, IsNil)
	c.Assert(ctx.IsAborted(), Equals, true)
	c.Assert(ctx.Writer.(*MockReponseWriter).Buffer.String(), Not(HasLen), 0)
	c.Assert(ctx.Writer.Header(), DeepEquals, http.Header{
		"Content-Type": []string{"text/html"},
	})
}

func (s *ProxySuite) TestDefaultHandler(c *C) {
	ctx := s.newMockContext("GET", "http://foo.gop.kg/bar!master")

	p := &Proxy{}

	err := p.Handle(ctx)
	c.Assert(err, IsNil)
	c.Assert(ctx.IsAborted(), Equals, false)
	c.Assert(ctx.Writer.(*MockReponseWriter).Buffer.String(), HasLen, 0)
}

func (s *ProxySuite) TestDoUploadPackInfoResponse(c *C) {
	ctx := s.newMockContext("POST", "http://tyba.gop.kg/git-fixture!master/info/refs")

	p := &Proxy{}

	err := p.Handle(ctx)
	c.Assert(err, IsNil)
	c.Assert(ctx.IsAborted(), Equals, true)
	c.Assert(ctx.Writer.(*MockReponseWriter).Buffer.String(), Not(HasLen), 0)
	c.Assert(ctx.Writer.Header(), DeepEquals, http.Header{
		"Content-Type": []string{"application/x-git-upload-pack-advertisement"},
	})
}

func (s *ProxySuite) TestDoUploadPackResponse(c *C) {
	ctx := s.newMockContext("POST", "http://tyba.gop.kg/git-fixture!master/git-upload-pack")

	p := &Proxy{}

	err := p.Handle(ctx)
	c.Assert(err, IsNil)
	c.Assert(ctx.IsAborted(), Equals, true)
	c.Assert(ctx.Writer.(*MockReponseWriter).Buffer.String(), Not(HasLen), 0)
	c.Assert(ctx.Writer.Header(), DeepEquals, http.Header{
		"Content-Type": []string{"application/x-git-upload-pack-result"},
	})
}

func (s *ProxySuite) TestDoUploadPackResponseNotFound(c *C) {
	ctx := s.newMockContext("POST", "http://qux.gop.kg/foo!baz/git-upload-pack")

	p := &Proxy{}

	err := p.Handle(ctx)
	c.Assert(err, IsNil)
	c.Assert(ctx.IsAborted(), Equals, true)
	c.Assert(ctx.Writer.(*MockReponseWriter).Buffer.String(), Not(HasLen), 0)
	c.Assert(ctx.Writer.Header(), DeepEquals, http.Header{
		"Content-Type":     []string{"text/plain; charset=utf-8"},
		"Www-Authenticate": []string{"Basic realm=\"GoPkg\""},
	})
}

func (s *ProxySuite) newMockContext(method, urlStr string) *gin.Context {
	ctx := &gin.Context{}
	ctx.Request, _ = http.NewRequest(method, urlStr, nil)
	ctx.Writer = &MockReponseWriter{}

	return ctx
}

type MockReponseWriter struct {
	StatusCode int
	Buffer     bytes.Buffer
	Headers    http.Header
}

func (w *MockReponseWriter) Header() http.Header {
	if len(w.Headers) == 0 {
		w.Headers = make(http.Header, 0)
	}

	return w.Headers
}

func (w *MockReponseWriter) Write(p []byte) (int, error) {
	return w.Buffer.Write(p)
}

func (w *MockReponseWriter) WriteString(s string) (int, error) {
	return w.Buffer.WriteString(s)
}

func (w *MockReponseWriter) WriteHeader(code int) {
	w.StatusCode = code
}

func (w *MockReponseWriter) CloseNotify() <-chan bool { return nil }
func (w *MockReponseWriter) Status() int              { return 0 }
func (w *MockReponseWriter) Size() int                { return 0 }
func (w *MockReponseWriter) Written() bool            { return false }
func (w *MockReponseWriter) WriteHeaderNow()          {}
func (w *MockReponseWriter) Flush()                   {}
func (w *MockReponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

package gopkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
}

func NewServer() *Server {
	return &Server{
		Engine: gin.Default(),
	}
}

func (s *Server) Run(addr ...string) error {
	s.Use(ProxyMiddleware())
	s.Any("/*default", func(c *gin.Context) {
		message := c.Param("default")
		c.String(http.StatusOK, message)
	})

	return s.Engine.Run(addr...)
}

func (s *Server) RunTLS(addr, certFile, keyFile string) error {
	return s.Engine.RunTLS(addr, certFile, keyFile)
}

func ProxyMiddleware() gin.HandlerFunc {
	proxy := &Proxy{}
	return func(c *gin.Context) {
		if err := proxy.Handle(c); err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.Next()
	}
}

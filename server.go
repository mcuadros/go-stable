package gopkg

import (
	"fmt"

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
	s.setTemplates()
	s.Use(ProxyMiddleware())
	s.NoRoute(s.Public)

	return s.Engine.Run(addr...)
}

func (s *Server) RunTLS(addr, certFile, keyFile string) error {
	s.setTemplates()
	s.Use(ProxyMiddleware())
	s.NoRoute(s.Public)

	return s.Engine.RunTLS(addr, certFile, keyFile)
}

func ProxyMiddleware() gin.HandlerFunc {
	proxy := &Proxy{}
	return func(c *gin.Context) {
		if err := proxy.Handle(c); err != nil {
			c.AbortWithError(500, err)
			return
		}

		fmt.Println(c.Get("package"))

		c.Next()
	}
}

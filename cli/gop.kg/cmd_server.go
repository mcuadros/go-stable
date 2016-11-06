package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mcuadros/gop.kg"
)

type ServerCommand struct {
	Addr     string `long:"addr" default:":8080" description:"http server addr"`
	CertFile string `long:"cert" description:"TLS certificate file path."`
	KeyFile  string `long:"key" description:"TLS key file path."`
	s        *gopkg.Server
}

func (c *ServerCommand) Execute(args []string) error {
	gin.SetMode(gin.ReleaseMode)

	c.s = gopkg.NewServer()
	//return c.s.Run(c.Addr)
	return c.s.RunTLS(c.Addr, c.CertFile, c.KeyFile)
}

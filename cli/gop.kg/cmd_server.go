package main

import (
	"time"

	"github.com/mcuadros/gop.kg"
)

type ServerCommand struct {
	Addr     string `long:"addr" default:":8080" description:"http server addr"`
	CertFile string `long:"cert" description:"TLS certificate file path."`
	KeyFile  string `long:"key" description:"TLS key file path."`
	s        *gopkg.Server
}

func (c *ServerCommand) Execute(args []string) error {
	c.s = gopkg.NewDefaultServer("example.com")
	c.s.Addr = c.Addr
	c.s.WriteTimeout = 15 * time.Second
	c.s.ReadTimeout = 15 * time.Second

	return c.s.ListenAndServeTLS(c.CertFile, c.KeyFile)
}

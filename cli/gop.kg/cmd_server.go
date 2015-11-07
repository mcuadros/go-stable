package main

import "github.com/mcuadros/gop.kg"

type ServerCommand struct {
	Addr string `long:"addr" default:":8080" description:"http server addr"`

	s *gopkg.Server
}

func (c *ServerCommand) Execute(args []string) error {
	c.s = gopkg.NewServer()
	return c.s.Run(c.Addr)
}

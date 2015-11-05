package main

import "github.com/mcuadros/gop.kg"

type ProxyCommand struct {
	Addr string `long:"addr" default:":8080" description:"http server addr"`

	p *gopkg.Proxy
}

func (c *ProxyCommand) Execute(args []string) error {
	c.p = gopkg.NewProxy()
	return c.p.Start(c.Addr)
}

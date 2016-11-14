package main

import (
	"crypto/tls"
	"fmt"

	"github.com/mcuadros/go-stable"

	"github.com/Sirupsen/logrus"
	"github.com/dkumor/acmewrapper"
	"github.com/urfave/negroni"
)

type ServerCommand struct {
	Host         string `long:"host" description:"host of the server"`
	Server       string `long:"server" default:"github.com" description:"repository git server"`
	Organization string `long:"organazation" description:"repository organization"`
	Repository   string `long:"repository" default:"github.com" description:"repository name"`

	Addr     string `long:"addr" default:":8080" description:"http server addr"`
	CertFile string `long:"cert" description:"TLS certificate file path."`
	KeyFile  string `long:"key" description:"TLS key file path."`

	LogLevel  string `long:"log-level" default:"info" description:"log level, values: debug, info, warn or panic"`
	LogFormat string `long:"log-format" default:"text" description:"log format, values: text or json"`

	s *stable.Server
}

func (c *ServerCommand) Execute(args []string) error {
	if err := c.buildServer(); err != nil {
		return err
	}

	if err := c.buildMiddleware(); err != nil {
		return err
	}

	return c.listen()
}

func (c *ServerCommand) buildServer() error {
	if c.Host == "" {
		return fmt.Errorf("missing host name, please set `--host`")
	}

	c.s = stable.NewDefaultServer(c.Host)
	c.s.Addr = c.Addr
	c.s.Default.Server = c.Server
	c.s.Default.Organization = c.Organization
	c.s.Default.Repository = c.Repository

	return nil
}

func (c *ServerCommand) buildMiddleware() error {
	logger, err := c.getLogrusMiddleware()
	if err != nil {
		return err
	}

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(logger)
	n.UseHandler(c.s.Handler)

	c.s.Handler = n

	return nil
}

func (c *ServerCommand) getLogrusMiddleware() (negroni.Handler, error) {
	level, err := c.getLogLevel()
	if err != nil {
		return nil, err
	}

	format, err := c.getLogFormat()
	if err != nil {
		return nil, err
	}

	return NewLogger(level, format), nil
}

func (c *ServerCommand) getLogLevel() (level logrus.Level, err error) {
	switch c.LogLevel {
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warn":
		level = logrus.WarnLevel
	case "panic":
		level = logrus.PanicLevel
	default:
		err = fmt.Errorf("invalid log-level, %q", c.LogLevel)
	}

	return
}

func (c *ServerCommand) getLogFormat() (format logrus.Formatter, err error) {
	switch c.LogFormat {
	case "text":
		format = &logrus.TextFormatter{}
	case "json":
		format = &logrus.JSONFormatter{}
	default:
		err = fmt.Errorf("invalid log-format, %q", c.LogLevel)
	}

	return
}

func (c *ServerCommand) listen() error {
	acme, err := c.getACME()
	if err != nil {
		return nil
	}

	c.s.TLSConfig = acme.TLSConfig()

	listener, err := tls.Listen("tcp", c.Addr, c.s.TLSConfig)
	if err != nil {
		return err
	}

	return c.s.Serve(listener)
}

func (c *ServerCommand) getACME() (*acmewrapper.AcmeWrapper, error) {
	return acmewrapper.New(acmewrapper.Config{
		Domains:     []string{c.Host},
		Address:     c.Addr,
		TLSCertFile: c.CertFile,
		TLSKeyFile:  c.KeyFile,
		TOSCallback: acmewrapper.TOSAgree,
	})
}

package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/dkumor/acmewrapper"
	"github.com/mcuadros/go-stable"
	"github.com/urfave/negroni"
)

const (
	BaseRoute       = "/{repository:[a-z0-9-/]+}.{version:v[0-9.]+}"
	BaseRouteOrg    = "/{org:[a-z0-9-]+}/{repository:[a-z0-9-/]+}.{version:v[0-9.]+}"
	BaseRouteSrvOrg = "/{srv:[a-z0-9-.]+}/{org:[a-z0-9-]+}/{repository:[a-z0-9-/]+}.{version:v[0-9.]+}"
)

type ServerCommand struct {
	Host         string `long:"host" description:"host of the server"`
	Server       string `long:"server" default:"github.com" description:"repository git server"`
	Organization string `long:"organization" description:"repository organization"`
	Repository   string `long:"repository" default:"github.com" description:"repository name"`
	BaseRoute    string `long:"base-route" description:"base gorilla/mux route"`

	Addr         string `long:"addr" default:":443" description:"http server addr"`
	RedirectAddr string `long:"redirect-addr" description:"http to https redirect server addr"`
	CertFolder   string `long:"certs" default:"/certificates" description:"TLS certificate folder"`

	LogLevel  string `long:"log-level" default:"info" description:"log level, values: debug, info, warn or panic"`
	LogFormat string `long:"log-format" default:"text" description:"log format, values: text or json"`

	s        *stable.Server
	redirect *http.Server
}

func (c *ServerCommand) Execute(args []string) error {
	if err := c.buildServer(); err != nil {
		return err
	}

	if err := c.buildMiddleware(); err != nil {
		return err
	}

	c.buildRedirectHTTP()
	return c.listen()
}

func (c *ServerCommand) buildServer() error {
	if c.Host == "" {
		return fmt.Errorf("missing host name, please set `--host`")
	}

	if c.BaseRoute == "" {
		c.BaseRoute = c.getBaseRoute()
	}

	c.s = stable.NewServer(c.BaseRoute, c.Host)
	c.s.Addr = c.Addr
	c.s.Default.Server = c.Server
	c.s.Default.Organization = c.Organization
	c.s.Default.Repository = c.Repository

	return nil
}

func (c *ServerCommand) getBaseRoute() string {
	if c.BaseRoute != "" {
		return c.BaseRoute
	}

	if c.Server == "" {
		return BaseRouteSrvOrg
	}

	if c.Organization == "" {
		return BaseRouteOrg
	}

	return BaseRoute
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
		return err
	}

	c.s.TLSConfig = acme.TLSConfig()

	listener, err := tls.Listen("tcp", c.Addr, c.s.TLSConfig)
	if err != nil {
		return err
	}

	go c.listenRedirectHTTP()
	return c.s.Serve(listener)
}

func (c *ServerCommand) getACME() (*acmewrapper.AcmeWrapper, error) {
	return acmewrapper.New(acmewrapper.Config{
		Domains:          []string{c.Host},
		Address:          c.Addr,
		TLSCertFile:      filepath.Join(c.CertFolder, "cert.pem"),
		TLSKeyFile:       filepath.Join(c.CertFolder, "key.pem"),
		RegistrationFile: filepath.Join(c.CertFolder, "user.reg"),
		PrivateKeyFile:   filepath.Join(c.CertFolder, "private.pem"),
		TOSCallback:      acmewrapper.TOSAgree,
	})
}

func (c *ServerCommand) listenRedirectHTTP() {
	if c.redirect == nil {
		return
	}

	c.redirect.ListenAndServe()
}

func (c *ServerCommand) buildRedirectHTTP() {
	if c.RedirectAddr == "" {
		return
	}

	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		url := req.URL
		url.Scheme = "https"
		url.Host = c.Host

		_, port, _ := net.SplitHostPort(c.Addr)
		if port != "443" {
			url.Host += ":" + port
		}

		http.Redirect(w, req, url.String(), http.StatusMovedPermanently)
	})

	c.redirect = &http.Server{
		Addr:    c.RedirectAddr,
		Handler: m,
	}
}

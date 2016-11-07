package main

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/negroni"
)

type Logger struct {
	Logger *logrus.Logger
}

func NewLogger(level logrus.Level, formatter logrus.Formatter) *Logger {
	log := logrus.New()
	log.Level = level
	log.Formatter = formatter

	return &Logger{Logger: log}
}

func (m *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request, n http.HandlerFunc) {
	remoteAddr := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		remoteAddr = realIP
	}

	f := logrus.Fields{}
	f["request"] = r.URL.Path
	f["method"] = r.Method
	f["remote"] = remoteAddr
	logrus.NewEntry(m.Logger).WithFields(f).Debug("new request accepted")

	start := time.Now()
	n(w, r)

	status := w.(negroni.ResponseWriter).Status()

	f["status"] = status
	f["elapsed"] = time.Since(start)

	entry := logrus.NewEntry(m.Logger).WithFields(f)
	if status != 200 {
		entry.Warning("completed handling request, with errors")
		return
	}

	entry.Info("completed handling request")
}

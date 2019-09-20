package ui

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
)

type handlerWrapper struct {
	handler http.Handler
}

var epoch = time.Unix(0, 0).Format(time.RFC1123)

func Serve(mux *http.ServeMux) *http.Server {
	l, e := net.Listen("tcp", "127.0.0.1:0")
	errs.Panic("Bind port", e)

	s := &http.Server{Handler: mux}
	s.Addr = "http://" + l.Addr().String()
	log.Println("Serving on:", s.Addr)

	go func() {
		e := s.Serve(l)
		errs.Log("Serve", e)
	}()

	return s
}

func (m *Manager) Handle(path string, handler func(http.ResponseWriter, *http.Request)) {
	m.mux.HandleFunc(path, wrapHandlerFunc(handler))
}

func (m *Manager) StripPrefix(path string, handler func(http.ResponseWriter, *http.Request)) {
	m.mux.Handle(path, wrapHandler(http.StripPrefix(path, http.HandlerFunc(handler))))
}

func (m *Manager) Call(path string, fn func(*Window)) {
	m.Handle(path, wrapHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(m.window)
	}))
}

func (h *handlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setNoCacheHeaders(w, r)
	h.handler.ServeHTTP(w, r)
}

func wrapHandler(h http.Handler) http.Handler {
	return &handlerWrapper{handler: h}
}

func wrapHandlerFunc(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setNoCacheHeaders(w, r)
		fn(w, r)
	}
}

func setNoCacheHeaders(w http.ResponseWriter, r *http.Request) {
	// remove header props
	w.Header().Del("ETag")
	w.Header().Del("If-Modified-Since")
	w.Header().Del("If-Match")
	w.Header().Del("If-None-Match")
	w.Header().Del("If-Unmodified-Since")

	// set header props
	w.Header().Set("Expires", epoch)
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")
}

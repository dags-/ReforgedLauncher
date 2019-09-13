package ui

import (
	"log"
	"net"
	"net/http"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
)

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
	m.mux.HandleFunc(path, handler)
}

func (m *Manager) StripPrefix(path string, handler func(http.ResponseWriter, *http.Request)) {
	m.mux.Handle(path, http.StripPrefix(path, http.HandlerFunc(handler)))
}

func (m *Manager) Call(path string, fn func(*Window)) {
	m.Handle(path, func(w http.ResponseWriter, r *http.Request) {
		fn(m.window)
	})
}

package ui

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	rice "github.com/GeertJohan/go.rice"
	"github.com/dags-/systray"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
)

type Manager struct {
	lock   *sync.RWMutex
	mux    *http.ServeMux
	srv    *http.Server
	name   string
	icon   string
	tray   string
	appDir string
	window *Window
}

func NewManager(appDir, name string, box *rice.Box) *Manager {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(box.HTTPBox()))
	return &Manager{
		lock:   &sync.RWMutex{},
		mux:    mux,
		srv:    Serve(mux),
		name:   name,
		appDir: appDir,
		icon:   filepath.Join(appDir, "Assets", icon("icon")),
		tray:   filepath.Join(appDir, "Assets", icon("tray")),
	}
}

func icon(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".ico"
	}
	return name + ".png"
}

func (m *Manager) Address() string {
	return m.srv.Addr
}

func (m *Manager) Window() *Window {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.window
}

func (m *Manager) HasWindow() bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.window == nil {
		return false
	}
	m.window.l.Lock()
	defer m.window.l.Unlock()
	return m.window.cmd != nil && m.window.events != nil
}

func (m *Manager) CloseWindow() {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.window != nil {
		m.window.Exit()
		m.window = nil
	}
}

func (m *Manager) Home() {
	e := m.Attach(&Settings{
		Url:       "/home",
		Width:     800,
		Height:    480,
		Resizable: true,
	})
	errs.Log("Attach window", e)
}

func (m *Manager) Progress(path string) {
	e := m.Attach(&Settings{
		Url:        "/progress#" + path,
		Width:      800,
		Height:     420,
		Borderless: true,
	})
	errs.Log("Attach window", e)
}

func (m *Manager) Attach(settings *Settings) error {
	m.CloseWindow()

	m.lock.Lock()
	defer m.lock.Unlock()

	settings.Title = m.name + settings.Title
	settings.Url = m.srv.Addr + settings.Url
	settings.Icon = m.icon

	if settings.Resizable {
		cfg := Load(m.appDir)
		settings.Width = cfg.WindowWidth
		settings.Height = cfg.WindowHeight
	}

	w, e := startWindow(settings)
	if e != nil {
		return e
	}

	m.window = w
	return nil
}

func (m *Manager) Exit() {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.window != nil {
		m.window.Exit()
		m.window = nil
	}
	systray.Quit()
	os.Exit(0)
}

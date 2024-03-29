package ui

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"

	rice "github.com/GeertJohan/go.rice"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/platform"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
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
	mux.Handle("/", wrapHandler(http.FileServer(box.HTTPBox())))
	return &Manager{
		lock:   &sync.RWMutex{},
		mux:    mux,
		srv:    Serve(mux),
		name:   name,
		appDir: appDir,
		icon:   filepath.Join(appDir, "Assets", platform.Icon("icon")),
		tray:   filepath.Join(appDir, "Assets", platform.Icon("tray")),
	}
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
		Height:     480,
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

	log.Println(settings.Url)

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
		log.Println("Closing window")
		terminate(m.window)
	}

	tasks.Shutdown()
}

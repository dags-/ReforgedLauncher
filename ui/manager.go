package ui

import (
	"net/http"
	"path/filepath"
	"sync"

	rice "github.com/GeertJohan/go.rice"
	"github.com/dags-/webview"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
)

type Manager struct {
	name   string
	appDir string
	mux    *http.ServeMux
	srv    *http.Server
	lock   *sync.RWMutex
	window *Window
	icon   []byte
	open   chan interface{}
	close  chan interface{}
	exit   chan interface{}
	create chan Settings
}

func NewManager(appdDir, name string, box *rice.Box) *Manager {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(box.HTTPBox()))
	return &Manager{
		mux:    mux,
		name:   name,
		appDir: appdDir,
		srv:    Serve(mux),
		icon:   getIcon(box),
		open:   make(chan interface{}),
		close:  make(chan interface{}),
		exit:   make(chan interface{}),
		create: make(chan Settings),
	}
}

func (m *Manager) Address() string {
	return m.srv.Addr
}

func (m *Manager) Home() {
	m.NewWindow(Settings{
		Path:      "/home",
		Width:     800,
		Height:    480,
		Resizable: true,
	})
}

func (m *Manager) Progress(path string) {
	m.NewWindow(Settings{
		Path:       "/progress#" + path,
		Width:      800,
		Height:     420,
		Borderless: true,
	})
}

func (m *Manager) NewSizedWindow(settings Settings) {
	cfg := Load(m.appDir)
	settings.Width = cfg.WindowWidth
	settings.Height = cfg.WindowHeight
	m.NewWindow(settings)
}

func (m *Manager) NewWindow(settings Settings) {
	go func() {
		m.create <- settings
	}()
}

func (m *Manager) Window() *Window {
	return m.window
}

func (m *Manager) Run() {
	defer close(m.open)
	defer close(m.exit)
	defer close(m.close)
	defer close(m.create)
	for {
		m.handleWindow()
		m.handleTray()
	}
}

func (m *Manager) Open() {
	ch := m.open
	if ch != nil {
		ch <- nil
	}
}

func (m *Manager) Close() {
	m.close <- nil
}

func (m *Manager) Exit() {
	m.exit <- nil
}

func (m *Manager) createWindow(settings Settings) *Window {
	w := &Window{
		w: webview.New(webview.Settings{
			Title:      m.name + settings.Title,
			URL:        m.srv.Addr + settings.Path,
			Width:      settings.Width,
			Height:     settings.Height,
			Resizable:  settings.Resizable,
			Borderless: settings.Borderless,
		}),
		b: true,
	}
	e := w.w.SetIcon(filepath.Join(m.appDir, "icon.ico"))
	errs.Log("Set Icon", e)
	return w
}

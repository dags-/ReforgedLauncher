package launcher

import (
	"github.com/GeertJohan/go.rice"

	"github.com/Conquest-Reforged/ReforgedLauncher/instance/repo"
	"github.com/Conquest-Reforged/ReforgedLauncher/ui"
)

type Launcher struct {
	*Properties
	r  *repo.Repository
	wm *ui.Manager
}

type Properties struct {
	AppDir      string `json:"-"`
	Branding    string `json:"branding"`
	ModPacksURL string `json:"modpacks_url"`
}

func NewLauncher(properties *Properties, box *rice.Box) *Launcher {
	r := repo.Open(properties.AppDir)
	wm := ui.NewManager(properties.AppDir, properties.Branding, box)
	return &Launcher{
		Properties: properties,
		r:          r,
		wm:         wm,
	}
}

func (l *Launcher) Run() {
	l.wm.Handle("/api/modpacks", l.modpacks)
	l.wm.Handle("/api/instances", l.instances)
	l.wm.Handle("/api/open/window", l.openWindow)
	l.wm.Handle("/api/open/folder", l.openFolder)
	l.wm.StripPrefix("/api/run/", l.run)
	l.wm.StripPrefix("/api/launch/", l.launch)
	l.wm.StripPrefix("/api/install/", l.install)
	l.wm.StripPrefix("/api/instance/", l.instance)

	c := ui.FirstLoad(l.AppDir, l.wm.Address())

	settings := ui.Settings{
		Path:      "/home",
		Width:     c.WindowWidth,
		Height:    c.WindowHeight,
		Resizable: true,
	}

	if c.AutoLaunch {
		last, e := l.LastInstance()
		if e == nil {
			settings.Path = "/progress#launch/" + last.Name
			settings.Width = 800
			settings.Height = 420
			settings.Resizable = false
			settings.Borderless = true
		}
	}

	l.wm.NewWindow(settings)
	l.wm.Run()
}

func (l *Launcher) Close() {
	l.wm.Close()
}

func (l *Launcher) Exit() {
	l.wm.Exit()
}

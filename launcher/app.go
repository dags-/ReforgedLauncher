package launcher

import (
	"github.com/Conquest-Reforged/ReforgedLauncher/instance/repo"
	"github.com/Conquest-Reforged/ReforgedLauncher/ui"
	"github.com/GeertJohan/go.rice"
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
	m := ui.NewManager(properties.AppDir, properties.Branding, box)
	return &Launcher{
		Properties: properties,
		r:          r,
		wm:         m,
	}
}

func (l *Launcher) Run() {
	l.wm.Handle("/api/modpacks", l.modpacks)
	l.wm.Handle("/api/instances", l.instances)
	l.wm.Handle("/api/folder/open", l.openFolder)
	l.wm.Handle("/api/window/open", l.openWindow)
	l.wm.Handle("/api/window/save", l.saveWindow)
	l.wm.Handle("/api/home/scale", l.scale)
	l.wm.StripPrefix("/api/run/", l.run)
	l.wm.StripPrefix("/api/launch/", l.launch)
	l.wm.StripPrefix("/api/install/", l.install)
	l.wm.StripPrefix("/api/instance/", l.instance)
	config := ui.FirstLoad(l.AppDir, l.wm.Address())
	l.launchWindow(config.AutoLaunch)
	l.wm.RunTray()
}

func (l *Launcher) Exit() {
	l.wm.Exit()
}

func (l *Launcher) Addr() string {
	return l.wm.Address()
}

func (l *Launcher) launchWindow(autoLaunch bool) {
	if autoLaunch {
		last, e := l.LastInstance()
		if e == nil {
			l.wm.Progress("launch/" + escape(last.Name))
			return
		}
	}
	l.wm.Home()
}

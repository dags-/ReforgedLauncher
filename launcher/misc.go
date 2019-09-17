package launcher

import (
	"net/http"

	"github.com/skratchdot/open-golang/open"

	"github.com/Conquest-Reforged/ReforgedLauncher/ui"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
)

func (l *Launcher) openWindow(w http.ResponseWriter, r *http.Request) {
	if l.wm.HasWindow() {
		return
	}
	_ = r.ParseForm()
	quick := r.FormValue("quick")
	if quick == "true" {
		l.launchWindow(true)
	} else {
		config := ui.Load(l.AppDir)
		l.launchWindow(config.AutoLaunch)
	}
}

func (l *Launcher) openFolder(w http.ResponseWriter, r *http.Request) {
	var install Install
	e := files.ParseJson(r.Body, &install)
	if e == nil {
		e = open.Start(install.Name)
		if e == nil {
			success(w, nil)
			return
		}
	}
	fail(w, e)
}

func (l *Launcher) saveWindow(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var cfg ui.Config
		e := files.ParseJson(r.Body, &cfg)
		if e == nil {
			config := ui.Load(l.AppDir)
			config.WindowWidth = cfg.WindowWidth
			config.WindowHeight = cfg.WindowHeight
			ui.Save(l.AppDir, config)
			success(w, "saved")
		} else {
			fail(w, e)
		}
	}
}

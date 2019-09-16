package launcher

import (
	"net/http"

	"github.com/skratchdot/open-golang/open"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
)

func (l *Launcher) openWindow(w http.ResponseWriter, r *http.Request) {
	l.wm.Home()
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

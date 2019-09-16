package launcher

import (
	"fmt"
	"net/http"

	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
)

type Install struct {
	Name string `json:"name"`
}

func (l *Launcher) install(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		l.getInstall(w, r)
	}
	if r.Method == http.MethodPost {
		l.postInstall(w, r)
	}
}

func (l *Launcher) getInstall(w http.ResponseWriter, r *http.Request) {
	listener := listener(l.wm.Window())

	listener.Stat("Loading instance", 0.33)
	instance, e := l.Instance(r.URL.Path)
	if e != nil {
		l.onError("Install", "Load instance", e)
		return
	}

	listener.Stat("Installing pack", 0.66)
	e = instance.InstallPack(l.r, listener)
	if e != nil {
		l.onError("Install", "Install modpack", e)
		return
	}

	listener.Stat("Saving instance", 1.0)
	e = l.SaveInstance(instance)
	if e != nil {
		l.onError("Install", "Save instance", e)
		return
	}

	listener.Wait()
	l.wm.Home()
}

func (l *Launcher) postInstall(w http.ResponseWriter, r *http.Request) {
	repo, e := modpack.ParseRepo(r.URL.Path)
	if e != nil {
		fail(w, e)
		return
	}

	remote, e := repo.Latest()
	if e != nil {
		fail(w, fmt.Errorf("get remote %s: %s", repo, e))
		return
	}

	var install Install
	e = files.ParseJson(r.Body, &install)
	if e != nil {
		fail(w, e)
		return
	}

	instances, e := l.LoadInstances()
	if e != nil {
		fail(w, fmt.Errorf("load instances: %s", e))
		return
	}

	install.Name = nameCleaner.ReplaceAllString(install.Name, "")
	if _, ok := instances[install.Name]; ok {
		fail(w, fmt.Errorf("instance already exists: %s", install.Name))
		return
	}

	instance := l.NewInstance(install.Name, remote)
	instances[instance.Name] = instance

	e = l.SaveInstances(instances)
	if e != nil {
		fail(w, fmt.Errorf("save instances %s: %s", install.Name, e))
		return
	}

	l.wm.Progress("install/" + instance.Name)
}

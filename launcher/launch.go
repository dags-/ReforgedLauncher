package launcher

import (
	"net/http"
	"time"

	"github.com/Conquest-Reforged/ReforgedLauncher/minecraft"
)

func (l *Launcher) run(w http.ResponseWriter, r *http.Request) {
	l.wm.Progress("launch/" + r.URL.Path)
}

func (l *Launcher) launch(w http.ResponseWriter, r *http.Request) {
	l.Launch(r.URL.Path)
}

func (l *Launcher) Launch(id string) {
	listener := listener(l.wm.Window())

	// load instance or create default
	listener.Stat("Loading instance", 0.2)
	instance, e := l.Instance(id)
	if e != nil {
		l.onError("Launch", "Load instance", e)
		return
	}

	listener.Stat("Checking for updates", 0.3)
	l.Update(instance)

	// make sure pack is downloaded
	listener.Stat("Installing modpack", 0.4)
	e = instance.InstallPack(l.r, listener)
	if e != nil {
		l.onError("Launch", "Install pack", e)
		return
	}

	// prepare installation
	listener.Stat("Preparing for launch", 0.6)
	installation := instance.Installation(l.r)
	e = instance.Prepare(installation, listener)
	if e != nil {
		l.onError("Launch", "Prepare launch", e)
		return
	}

	// timestamp latest used pack
	instance.LastUsed = time.Now()
	e = l.SaveInstance(instance)
	if e != nil {
		l.onError("Launch", "Save instance", e)
		return
	}

	// get mojang launcher
	launcher, e := minecraft.Launcher(l.AppDir)
	if e != nil {
		listener.Stat("Installing minecraft launcher", 0.9)
		launcher, e = minecraft.Install(l.AppDir, listener)
		if e != nil {
			l.onError("Launch", "Install mojang launcher", e)
			return
		}
	}

	listener.Stat("Launching", 1.0)
	listener.Wait()

	// run mojang launcher
	e = launcher.Launch(installation).Start()
	if e != nil {
		l.onError("Launch", "Launch instance", e)
		return
	}

	// close window
	l.wm.CloseWindow()
}

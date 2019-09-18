package launcher

import (
	"log"
	"net/http"
	"time"

	"github.com/Conquest-Reforged/ReforgedLauncher/minecraft"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
)

func (l *Launcher) run(w http.ResponseWriter, r *http.Request) {
	l.wm.Progress("launch/" + escape(r.URL.Path))
}

func (l *Launcher) launch(w http.ResponseWriter, r *http.Request) {
	log.Println("launch", r.URL.Path)
	l.Launch(r.URL.Path)
}

func (l *Launcher) Launch(id string) {
	id = unescape(id)
	listener := listener(l.wm.Window())

	// load instance or create default
	listener.Stat("Loading instance", 0.2)
	instance, e := l.Instance(id)
	if e != nil {
		l.onError("Command", "Load instance", e)
		return
	}

	if instance.AutoUpdate {
		listener.Stat("Checking for updates", 0.3)
		l.Update(instance)
	}

	// make sure pack is downloaded
	listener.Stat("Installing modpack", 0.4)
	e = instance.InstallPack(l.r, listener)
	if e != nil {
		l.onError("Command", "Install pack", e)
		return
	}

	// prepare installation
	listener.Stat("Preparing for launch", 0.6)
	installation := instance.Installation(l.r)
	e = instance.Prepare(installation, listener)
	if e != nil {
		l.onError("Command", "Prepare launch", e)
		return
	}

	// timestamp latest used pack
	instance.LastUsed = time.Now()
	e = l.SaveInstance(instance)
	if e != nil {
		l.onError("Command", "Save instance", e)
		return
	}

	// init minecraft
	mc, e := minecraft.Get(l.AppDir)
	if e != nil {
		l.onError("Minecraft", "Load minecraft", e)
		return
	}

	// get mojang launcher
	launcher, e := mc.GetLauncher(listener)
	if e != nil {
		l.onError("Command", "Install mojang launcher", e)
		return
	}

	listener.Stat("Launching", 1.0)
	listener.Wait()

	// run mojang launcher
	e = tasks.Start(launcher.Command(installation))
	if e != nil {
		l.onError("Command", "Command instance", e)
		return
	}

	// close window
	l.wm.CloseWindow()
}

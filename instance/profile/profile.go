package profile

import (
	"time"

	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/node"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
)

const (
	javaArgs = "-Xmx4G -XX:+UnlockExperimentalVMOptions -XX:+UseG1GC -XX:G1NewSizePercent=20" +
		" -XX:G1ReservePercent=20 -XX:MaxGCPauseMillis=50 -XX:G1HeapRegionSize=16M"
)

func Init(i *modpack.Installation) {
	out := files.MustFile(i.GameDir, "launcher_profiles.json")
	if files.Exists(out) {
		return
	}

	in := files.MustFile(i.AppDir, "Repository", "Bin", "launcher_profiles.json")
	e := tasks.CopyPath(in, out, nil)
	errs.Panic("Write profiles", e)
}

func Set(brand string, i *modpack.Installation) {
	path := files.MustFile(i.GameDir, "launcher_profiles.json")
	profile := node.FromJsonAt(path)
	if profile.Empty() {
		return
	}

	forge := profile.Get("profiles").Get("forge")
	if profile.Empty() {
		return
	}

	// Mon Jan 2 15:04:05 -0700 MST 2006
	stamp := time.Now().Format("2006-01-02T15:04:05.000Z")
	forge.Set("name", brand)
	forge.Set("lastUsed", stamp)
	forge.Set("javaArgs", javaArgs)

	e := profile.ToJsonAt(path)
	errs.Log("Save profile", e)
}

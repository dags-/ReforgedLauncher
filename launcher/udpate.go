package launcher

import (
	"log"
	"path/filepath"

	"github.com/Conquest-Reforged/ReforgedLauncher/instance"
	"github.com/Conquest-Reforged/ReforgedLauncher/instance/forge"
	"github.com/Conquest-Reforged/ReforgedLauncher/instance/pack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/version"
)

func (l *Launcher) Update(inst *instance.Instance) {
	remote, e := inst.ModPack.Repo.Latest()
	if e != nil {
		l.onError("Update", "Get latest release", e)
		return
	}

	index := inst.ModPack.Version.UpgradeIndex(remote.Version)
	if index == version.NoUpgrade {
		return
	}

	installation := inst.Installation(l.r)
	base, e := pack.Load(installation)
	if e != nil {
		l.onError("Update", "Load pack", e)
		return
	}

	if index < version.MajorUpgrade {
		log.Println("uninstalling pack")
		pack.Uninstall(installation, base)

		log.Println("uninstalling forge")
		forge.Uninstall(installation)

		instanceDir := filepath.Join(inst.GameDir, inst.ModPack.Version.String())
		_, e = files.RenameRel(instanceDir, remote.Version.String())
		if e != nil {
			l.onError("Update", "Rename dir", e)
			return
		}
	}

	inst.ModPack = remote.AsModPack()
	e = l.SaveInstance(inst)
	if e != nil {
		l.onError("Update", "Save instance", e)
		return
	}
}

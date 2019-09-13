package instance

import (
	"time"

	"github.com/Conquest-Reforged/ReforgedLauncher/instance/forge"
	"github.com/Conquest-Reforged/ReforgedLauncher/instance/pack"
	"github.com/Conquest-Reforged/ReforgedLauncher/instance/profile"
	"github.com/Conquest-Reforged/ReforgedLauncher/instance/repo"
	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

type Meta struct {
	*modpack.ModPack
	Name   string `json:"name"`
	Pack   string `json:"pack"`
	Cover1 string `json:"cover1"`
	Cover2 string `json:"cover2"`
}

type Instance struct {
	Name     string           `json:"name"`
	Image    string           `json:"image"`
	GameDir  string           `json:"game_dir"`
	LastUsed time.Time        `json:"last_used"`
	ModPack  *modpack.ModPack `json:"mod_pack"`
	Options  *pack.Options    `json:"options"`
	AppDir   string           `json:"-"`
}

func (i *Instance) Installation(repository *repo.Repository) *modpack.Installation {
	return &modpack.Installation{
		AppDir:   i.AppDir,
		GameDir:  files.MustDir(i.gameDir(), i.ModPack.Version.String()),
		PackDir:  repository.PackDir(i.ModPack),
		ForgeDir: repository.LoaderDir(i.ModPack),
	}
}

func (i *Instance) InstallPack(r *repo.Repository, listener progress.Listener) error {
	if !r.Has(i.ModPack) {
		e := r.Pull(i.ModPack.Repo, i.ModPack.Version, listener)
		if e != nil {
			return e
		}
	}

	base, e := pack.Load(i.Installation(r))
	if e != nil {
		return e
	}

	i.Options = base.Sync(i.Options)
	return nil
}

func (i *Instance) Prepare(installation *modpack.Installation, listener progress.Listener) error {
	listener.GlobalStatus("Loading base modpack")
	base, e := pack.Load(installation)
	if e != nil {
		return e
	}

	listener.GlobalStatus("Synchronizing user options")
	i.Options = base.Sync(i.Options)
	e = pack.Install(installation, base, listener)
	if e != nil {
		return e
	}

	if !forge.IsInstalled(installation) {
		listener.GlobalStatus("Installing forge")
		e = forge.Install(installation, listener)
		if e != nil {
			return e
		}
		listener.TaskStatus("")
	}

	listener.GlobalStatus("Setting launcher profile")
	profile.Set(i.ModPack.Repo.Name, installation)
	return nil
}

func (i *Instance) gameDir() string {
	if i.GameDir == "" {
		return files.MustDir(i.AppDir, "ModPacks", i.ModPack.Repo.Name)
	}
	return i.GameDir
}

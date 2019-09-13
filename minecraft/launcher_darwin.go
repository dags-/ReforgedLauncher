package minecraft

import (
	"os/exec"

	"github.com/Conquest-Reforged/ReforgedLauncher/instance/repo"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

func executable() string {
	return "MinecraftLauncher.app"
}

func (m *Meta) platform() *AppLink {
	return m.OSX
}

func install(file string, listener progress.Listener) (string, error) {
	out := tasks.RelDir(file, executable())
	progress.Status("Extracting launcher")
	return out, tasks.Unzip(file, out, listener)
}

func launch(path string, r *repo.Repo) *exec.Cmd {
	return exec.Command("open", path, "--workDir", r.GameDir)
}

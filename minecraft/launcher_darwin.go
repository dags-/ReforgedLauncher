package minecraft

import (
	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
	"os/exec"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

func executable() string {
	return "MinecraftLauncher.app"
}

func (m *Meta) platform() *AppLink {
	return m.OSX
}

func install(file string, listener progress.Listener) (string, error) {
	out := files.RelDir(file, executable())
	listener.GlobalStatus("Extracting launcher")
	return out, tasks.Unzip(file, out, listener)
}

func launch(path string, i *modpack.Installation) *exec.Cmd {
	return exec.Command("open", path, "--workDir", i.GameDir)
}

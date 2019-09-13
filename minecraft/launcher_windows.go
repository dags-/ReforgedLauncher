package minecraft

import (
	"os/exec"

	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

func executable() string {
	return "MinecraftLauncher.exe"
}

func (m *Meta) platform() *AppLink {
	return m.Windows
}

func install(path string, listener progress.Listener) (string, error) {
	return files.RenameRel(path, executable())
}

func launch(path string, i *modpack.Installation) *exec.Cmd {
	return exec.Command(path, "--workDir", i.GameDir)
}

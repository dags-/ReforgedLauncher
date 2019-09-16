package platform

import (
	"os/exec"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

var platform = &linux{}

type linux struct {
}

func (l linux) AppDir(name string) string {
	return files.MustDir(Home(), "."+name)
}

func (l linux) Icon(name string) string {
	return name + ".png"
}

func (l linux) HideConsole(cmd *exec.Cmd) {

}

func (l linux) LauncherName() string {
	return "MinecraftLauncher"
}

func (l linux) ExtractLauncher(path string, listener progress.Listener) (string, error) {
	return files.RenameRel(path, platform.LauncherName())
}

func (l linux) LaunchCmd(exe, workDir string) *exec.Cmd {
	return exec.Command("./"+exe, "--workDir", workDir)
}

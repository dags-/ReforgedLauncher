package platform

import (
	"os/exec"

	"github.com/Conquest-Reforged/ReforgedLauncher/minecraft"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
)

var platform = &darwin{}

type darwin struct {
}

func (d darwin) Icon(name string) string {
	return name + ".png"
}

func (d darwin) HideConsole(cmd *exec.Cmd) {

}

func (d darwin) AppLink(meta *minecraft.Meta) *minecraft.AppLink {
	return meta.OSX
}

func (d darwin) LauncherName() string {
	return "MinecraftLauncher.app"
}

func (d darwin) ExtractLauncher(path string, listener progress.Listener) (string, error) {
	out := files.RelDir(path, platform.LauncherName())
	listener.GlobalStatus("Extracting launcher")
	return out, tasks.Unzip(path, out, listener)
}

func (d darwin) LaunchCmd(exe, workDir string) *exec.Cmd {
	return exec.Command(exe, "--workDir", workDir)
}

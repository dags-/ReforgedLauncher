package platform

import (
	"os/exec"
	"syscall"

	"github.com/Conquest-Reforged/ReforgedLauncher/minecraft"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

var platform = &windows{}

type windows struct {
}

func (w windows) Icon(name string) string {
	return name + ".ico"
}

func (w windows) HideConsole(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

func (w windows) AppLink(meta *minecraft.Meta) *minecraft.AppLink {
	return meta.Windows
}

func (w windows) LauncherName() string {
	return "MinecraftLauncher.exe"
}

func (w windows) ExtractLauncher(path string, listener progress.Listener) (string, error) {
	return files.RenameRel(path, platform.LauncherName())
}

func (w windows) LaunchCmd(exe, workDir string) *exec.Cmd {
	return exec.Command("open", exe, "--workDir", workDir)
}

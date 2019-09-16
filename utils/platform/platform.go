package platform

import (
	"os/exec"

	"github.com/Conquest-Reforged/ReforgedLauncher/minecraft"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

type Platform interface {
	Icon(name string) string

	HideConsole(cmd *exec.Cmd)

	AppLink(meta *minecraft.Meta) *minecraft.AppLink

	LauncherName() string

	ExtractLauncher(path string, listener progress.Listener) (string, error)

	LaunchCmd(exe, workDir string) *exec.Cmd
}

func Icon(name string) string {
	return platform.Icon(name)
}

func HideConsole(cmd *exec.Cmd) {
	platform.HideConsole(cmd)
}

func AppLink(meta *minecraft.Meta) *minecraft.AppLink {
	return platform.AppLink(meta)
}

func LauncherName() string {
	return platform.LauncherName()
}

func ExtractLauncher(archive string, listener progress.Listener) (string, error) {

}

func LaunchCmd(exe, workDir string) *exec.Cmd {
	return platform.LaunchCmd(exe, workDir)
}

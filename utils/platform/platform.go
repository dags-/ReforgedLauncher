package platform

import (
	"os/exec"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

type Platform interface {
	Icon(name string) string

	HideConsole(cmd *exec.Cmd)

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

func LauncherName() string {
	return platform.LauncherName()
}

func ExtractLauncher(path string, listener progress.Listener) (string, error) {
	return platform.ExtractLauncher(path, listener)
}

func LaunchCmd(exe, workDir string) *exec.Cmd {
	return platform.LaunchCmd(exe, workDir)
}

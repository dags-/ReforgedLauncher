package platform

import (
	"os"
	"os/exec"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

type Platform interface {
	AppDir(name string) string

	Icon(name string) string

	HideConsole(cmd *exec.Cmd)

	LauncherName() string

	ExtractLauncher(path string, listener progress.Listener) (string, error)

	LaunchCmd(exe, workDir string) *exec.Cmd
}

func AppDir(name string) string {
	return platform.AppDir(name)
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

func Home() string {
	h, e := os.UserHomeDir()
	errs.Panic("User Home Dir", e)
	return h
}

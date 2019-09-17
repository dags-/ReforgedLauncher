package platform

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
)

var platform = &darwin{}

type darwin struct {
}

func (d darwin) AppDir(name string) string {
	return files.MustDir(Home(), "Library", "Application Support", name)
}

func (d darwin) Icon(name string) string {
	return name + ".png"
}

func (d darwin) HideConsole(cmd *exec.Cmd) {

}

func (d darwin) LauncherName() string {
	return "MinecraftLauncher.app"
}

func (d darwin) ExtractLauncher(path string, listener progress.Listener) (string, error) {
	out := files.RelDir(path, platform.LauncherName())
	listener.GlobalStatus("Extracting launcher")
	e := tasks.Unzip(path, files.MustDir(out, "Contents"), listener)
	files.Del(path)
	if e != nil {
		return "", e
	}

	exe := filepath.Join(out, "Contents", "MacOS", "launcher")
	e = os.Chmod(exe, os.ModePerm)
	if e != nil {
		return "", nil
	}

	return out, e
}

func (d darwin) RunExecutable(exe string, args ...string) *exec.Cmd {
	cmd := make([]string, len(args)+2)
	cmd[0] = exe
	cmd[1] = "--args"
	for i, a := range args {
		cmd[2+i] = a
	}
	return exec.Command("open", cmd...)
}

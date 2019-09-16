package minecraft

import (
	"fmt"
	"os/exec"

	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/platform"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

type MojangLauncher struct {
	path string
}

func Launcher(appDir string) (*MojangLauncher, error) {
	path := files.MustFile(appDir, "Launcher", platform.LauncherName())
	if !files.Exists(path) {
		return nil, fmt.Errorf("mojang launcher not found")
	}
	return &MojangLauncher{path: path}, nil
}

func Install(appDir string, listener progress.Listener) (*MojangLauncher, error) {
	file, e := download(appDir, listener)
	if e != nil {
		return nil, e
	}

	file, e = platform.ExtractLauncher(file, listener)
	if e != nil {
		return nil, e
	}

	return &MojangLauncher{path: file}, nil
}

func (l *MojangLauncher) Launch(i *modpack.Installation) *exec.Cmd {
	cmd := platform.RunExecutable(l.path, "--workDir", i.GameDir)
	cmd.Dir = i.AppDir
	return cmd
}

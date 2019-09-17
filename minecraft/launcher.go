package minecraft

import (
	"errors"
	"os/exec"

	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/platform"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
)

type Launcher struct {
	path string
}

func (l *Launcher) Command(i *modpack.Installation) *exec.Cmd {
	cmd := platform.RunExecutable(l.path, "--workDir", i.GameDir)
	cmd.Dir = i.AppDir
	return cmd
}

func installLauncher(appDir string, listener progress.Listener) (string, error) {
	file, e := downloadLauncher(appDir, listener)
	if e != nil {
		return "", e
	}

	file, e = platform.ExtractLauncher(file, listener)
	if e != nil {
		return "", e
	}

	return file, nil
}

func downloadLauncher(appDir string, listener progress.Listener) (string, error) {
	m, e := getMeta()
	if e != nil {
		return "", e
	}

	link := m.getAppLink()
	path := files.TempFile(appDir, "Launcher")

	listener.GlobalStatus("Downloading Minecraft launcher")
	e = tasks.Download(link.AppLink, path, listener)
	if e != nil {
		return "", e
	}

	listener.GlobalStatus("Checking hash")
	if !files.CheckMD5(path, link.DownloadHash) {
		return "", errors.New("invalid md5 hash")
	}

	return path, nil
}

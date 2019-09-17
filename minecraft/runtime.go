package minecraft

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
)

type Runtime struct {
	path string
}

func (r *Runtime) Command(args ...string) *exec.Cmd {
	return exec.Command(r.path, args...)
}

func installRuntime(appdir string, meta *AppMeta, listener progress.Listener) (string, error) {
	listener.GlobalStatus("Installing java runtime")
	dir := files.MustDir(appdir, "Launcher", "runtime")
	lzma, e := downloadRuntime(meta, dir, listener)
	if e != nil {
		return "", e
	}

	jre, e := extractRuntime(lzma, dir, meta, listener)
	if e != nil {
		return "", e
	}

	java := filepath.Join(jre, "bin", exe())
	if !files.Exists(java) {
		return "", fmt.Errorf("java not found at" + java)
	}

	_ = os.Chmod(java, os.ModePerm)

	return java, nil
}

func downloadRuntime(meta *AppMeta, dir string, listener progress.Listener) (string, error) {
	listener.TaskStatus("Downloading runtime")
	lzma := files.TempFile(dir)
	e := tasks.Download(meta.X64.JRE.Url, lzma, listener)
	if e != nil {
		return "", e
	}

	listener.TaskStatus("Checking runtime")
	if !files.CheckMD5(lzma, meta.X64.JRE.Sha1) {
		return "", fmt.Errorf("mismatching sha1 checksum")
	}

	return lzma, nil
}

func extractRuntime(path, dir string, meta *AppMeta, listener progress.Listener) (string, error) {
	zip := files.TempFile(dir)
	defer files.Del(path)
	defer files.Del(zip)

	listener.TaskStatus("Decompressing runtime")
	e := tasks.Decompress(path, zip, listener)
	if e != nil {
		return "", e
	}

	listener.TaskStatus("Unzipping runtime")
	jre := files.MustDir(dir, "jre-x64", meta.X64.JRE.Version)
	e = tasks.Unzip(zip, jre, listener)
	if e != nil {
		return "", e
	}

	return jre, nil
}

func exe() string {
	if runtime.GOOS == "windows" {
		return "javaw.exe"
	}
	return "javaw"
}

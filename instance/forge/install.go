package forge

import (
	"bufio"
	"fmt"
	"html"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Conquest-Reforged/ReforgedLauncher/instance/profile"
	"github.com/Conquest-Reforged/ReforgedLauncher/minecraft"
	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/platform"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

func Install(i *modpack.Installation, listener progress.Listener) error {
	dir := filepath.Join(i.AppDir, "Repository", "Bin")
	matches := files.Match(dir, "ForgeWrapper.*?\\.jar")
	if len(matches) == 0 {
		return fmt.Errorf("forge wrapper not found")
	}

	wrapper := matches[len(matches)-1]
	installer, ok := findForgeInstaller(i)
	if !ok {
		return fmt.Errorf("forge installer not found")
	}

	e := installForge(i, wrapper, installer, listener)
	if e != nil {
		return e
	}
	return nil
}

func installForge(i *modpack.Installation, wrapper, installer string, listener progress.Listener) error {
	mc, e := minecraft.Get(i.AppDir)
	if e != nil {
		return e
	}

	r, e := mc.GetRuntime(listener)
	if e != nil {
		return fmt.Errorf("java not installed %s", e)
	}

	// negative number == indeterminate progress bar
	listener.TaskProgress(-1)

	// add the wrapper and installer jars to the classpath
	classpath := buildClassPath(wrapper, installer)
	cmd := buildCommand(r, classpath, i.GameDir)
	drain(cmd, listener)

	// forge installer needs the launcher_profiles.json to be present
	profile.Init(i)

	// run forge installer
	return cmd.Run()
}

func buildCommand(java *minecraft.Runtime, classpath, gameDir string) *exec.Cmd {
	cmd := java.Command("-classpath", classpath, "Main", gameDir)
	platform.HideConsole(cmd)
	return cmd
}

func buildClassPath(wrapper, installer string) string {
	if runtime.GOOS == "windows" {
		return wrapper + ";" + installer
	}
	return wrapper + ":" + installer
}

func findForgeInstaller(i *modpack.Installation) (string, bool) {
	matches := files.Match(i.ForgeDir, ".*\\.jar$")
	if len(matches) == 0 {
		return "", false
	}
	return matches[0], true
}

func drain(cmd *exec.Cmd, listener progress.Listener) {
	out, e := cmd.StdoutPipe()
	errs.Log("StdoutPipe", e)
	if e == nil {
		go drainOutput(out, listener)
	}

	err, e := cmd.StderrPipe()
	errs.Log("StderrPipe", e)
	if e == nil {
		go drainOutput(err, listener)
	}
}

func drainOutput(r io.ReadCloser, listener progress.Listener) {
	defer files.Close(r)
	s := bufio.NewScanner(r)
	for s.Scan() {
		text := s.Text()
		text = strings.TrimSpace(text)
		text = html.EscapeString(text)
		listener.TaskStatus(text)
	}
}

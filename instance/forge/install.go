package forge

import (
	"bufio"
	"fmt"
	"github.com/Conquest-Reforged/ReforgedLauncher/instance/profile"
	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
	"os/exec"
	"path/filepath"
	"strings"
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
	// add the wrapper and installer jars to the classpath
	classpath := fmt.Sprintf(`%s;%s`, wrapper, installer)
	cmd := buildCommand(classpath, i.GameDir)
	out, e := cmd.StdoutPipe()
	if e != nil {
		return e
	}

	// forge installer needs the launcher_profiles.json to be present
	profile.Init(i)

	// launch the forge installer
	e = cmd.Start()
	if e != nil {
		return e
	}

	// read the process output and feedback into this app
	s := bufio.NewScanner(out)
	listener.TaskProgress(-1)
	for s.Scan() {
		listener.TaskStatus(strings.TrimSpace(s.Text()))
	}

	return nil
}

func buildCommand(classpath, gameDir string) *exec.Cmd {
	return exec.Command("java", "-classpath", classpath, "Main", gameDir)
}

func findForgeInstaller(i *modpack.Installation) (string, bool) {
	matches := files.Match(i.ForgeDir, ".*\\.jar$")
	if len(matches) == 0 {
		return "", false
	}
	return matches[0], true
}

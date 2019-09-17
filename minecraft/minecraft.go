package minecraft

import (
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/platform"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

type Minecraft struct {
	appDir string
	meta   *AppMeta
}

var (
	lock = &sync.RWMutex{}
	mc   *Minecraft
)

func Get(appDir string) (*Minecraft, error) {
	lock.Lock()
	defer lock.Unlock()
	if mc == nil {
		m, e := newMinecraft(appDir)
		if e != nil {
			return nil, e
		}
		mc = m
	}
	return mc, nil
}

func newMinecraft(appDir string) (*Minecraft, error) {
	meta, e := getMeta()
	if e != nil {
		return nil, e
	}
	return &Minecraft{meta: meta.getAppLink(), appDir: appDir}, nil
}

func (m *Minecraft) GetLauncher(listener progress.Listener) (*Launcher, error) {
	path := filepath.Join(m.appDir, "Launcher", platform.LauncherName())
	if !files.Exists(path) {
		_, e := installLauncher(m.appDir, listener)
		if e != nil {
			return nil, e
		}
	}
	return &Launcher{path: path}, nil
}

func (m *Minecraft) GetRuntime(listener progress.Listener) (*Runtime, error) {
	java, e := exec.LookPath("java")
	if e != nil {
		return nil, e
	}

	if true {
		return &Runtime{path: java}, nil
	}

	path := filepath.Join(m.appDir, "Launcher", "runtime", "jre-x64", m.meta.X64.JRE.Version, "bin", exe())
	if !files.Exists(path) {
		_, e := installRuntime(m.appDir, m.meta, listener)
		if e != nil {
			return nil, e
		}
	}
	return &Runtime{path: path}, nil
}

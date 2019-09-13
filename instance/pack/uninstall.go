package pack

import (
	"path/filepath"

	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
)

func Uninstall(i *modpack.Installation, p *Pack) {
	del(i.GameDir, p.Required)
	for _, opt := range *p.Options {
		del(i.GameDir, opt.Files)
	}
}

func del(dir string, paths []*Path) {
	for _, p := range paths {
		path := filepath.Join(dir, p.Dest)
		files.Del(path)
	}
}

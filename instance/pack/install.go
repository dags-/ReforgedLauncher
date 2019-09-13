package pack

import (
	"os"
	"path/filepath"

	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
)

type Installer struct {
	add []*Path
	del []*Path
}

func Install(i *modpack.Installation, p *Pack, listener progress.Listener) error {
	// get lists of files that need either copying to the gameDir or removing from the gameDir
	add, del := getPaths(p, i.GameDir, i.PackDir)

	ch := make(chan float64)
	defer close(ch)
	go progress.Listen(ch, listener)

	a := len(add)
	d := len(del)
	count := float64(0)
	total := float64(a + d)

	// remove installed files that are no-longer wanted in the gameDir
	for _, f := range del {
		listener.TaskStatus("Deleting " + f.Dest)
		e := os.Remove(f.Dest)
		if e != nil {
			return e
		}
		count++
		ch <- count / total
	}

	// add files that are now required but don't currently exist in the gameDir
	for _, f := range add {
		listener.TaskStatus("Copying " + f.Dest)
		e := tasks.CopyPath(f.Src, f.Dest, listener)
		if e != nil {
			return e
		}
		count++
		ch <- count / total
	}

	return nil
}

func getPaths(p *Pack, gameDir, packDir string) ([]*Path, []*Path) {
	var add []*Path
	var del []*Path

	for _, f := range p.Required {
		src := filepath.Join(packDir, f.Src)
		dst := filepath.Join(gameDir, f.Dest)
		if files.Exists(dst) {
			continue
		}
		add = append(add, &Path{
			Src:  src,
			Dest: dst,
		})
	}

	for _, a := range *p.Options {
		enabled := isEnabled(p, a)
		for _, f := range a.Files {
			dst := filepath.Join(gameDir, f.Dest)
			exists := files.Exists(dst)

			if exists && !enabled {
				src := filepath.Join(packDir, f.Src)
				del = append(del, &Path{
					Src:  src,
					Dest: dst,
				})
			}

			if !exists && enabled {
				src := filepath.Join(packDir, f.Src)
				add = append(add, &Path{
					Src:  src,
					Dest: dst,
				})
			}
		}
	}

	return add, del
}

func isEnabled(p *Pack, a *Item) bool {
	if !a.Enabled {
		return false
	}
	for _, dep := range a.Dependencies {
		d, ok := (*p.Options)[dep]
		if !ok {
			return false
		}
		if !isEnabled(p, d) {
			return false
		}
	}
	return true
}

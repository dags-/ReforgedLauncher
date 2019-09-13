package pack

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
)

func Load(i *modpack.Installation) (*Pack, error) {
	t := &Pack{Required: []*Path{}, Options: &Options{}}
	e := readDir(t, nil, i.PackDir, i.PackDir)
	return t, e
}

func item(name string) *Item {
	return &Item{
		Name:         name,
		Enabled:      true,
		Files:        []*Path{},
		Dependencies: []string{},
	}
}

func (p *Pack) add(item *Item) {
	i, ok := (*p.Options)[item.Name]
	if !ok {
		(*p.Options)[item.Name] = item
		return
	}

	if item.Files != nil {
		i.Files = append(i.Files, item.Files...)
	}

	if item.Dependencies != nil {
		i.Dependencies = append(i.Dependencies, item.Dependencies...)
	}
}

func readDir(tree *Pack, parent *Item, root, dir string) error {
	fs, e := ioutil.ReadDir(dir)
	if e != nil {
		return e
	}

	for _, f := range fs {
		if f.IsDir() {
			e := addDir(tree, parent, root, dir, f.Name())
			if e != nil {
				return e
			}
		} else if f.Name() == "description.txt" {
			if parent != nil {
				path := filepath.Join(dir, f.Name())
				text, _ := files.ReadText(path)
				parent.Description = text
			}
		} else if strings.HasSuffix(f.Name(), ".dl.txt") {

		} else {
			addFile(tree, parent, root, dir, f.Name())
		}
	}
	return nil
}

func addDir(tree *Pack, parent *Item, root, dir, name string) error {
	if strings.HasPrefix(name, "_") {
		child := item(name[1:])
		e := readDir(tree, child, root, filepath.Join(dir, name))
		if e != nil {
			return e
		}
		if parent != nil {
			child.Dependencies = append(child.Dependencies, parent.Name)
		}
		tree.add(child)
		return nil
	} else {
		return readDir(tree, parent, root, filepath.Join(dir, name))
	}
}

func addFile(tree *Pack, parent *Item, root, dir, name string) {
	src := rel(root, filepath.Join(dir, name))
	path := &Path{
		Src:  src,
		Dest: dest(src),
	}
	if parent == nil {
		tree.Required = append(tree.Required, path)
	} else {
		parent.Files = append(parent.Files, path)
	}
}

func dest(path string) string {
	var dst []string
	parts := strings.Split(path, string(filepath.Separator))
	for _, part := range parts {
		if strings.HasPrefix(part, "_") {
			continue
		}
		dst = append(dst, part)
	}
	return filepath.Join(dst...)
}

func rel(dir, path string) string {
	r, e := filepath.Rel(dir, path)
	if e != nil {
		return path
	}
	return r
}

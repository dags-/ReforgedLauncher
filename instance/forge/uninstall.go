package forge

import (
	"archive/zip"
	"path/filepath"

	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/node"
)

func Uninstall(i *modpack.Installation) {
	version, ok := getForgeVersion(i)
	if !ok {
		return
	}

	dir := filepath.Join(i.GameDir, "versions", version)
	files.Del(dir)
}

func getForgeVersion(i *modpack.Installation) (string, bool) {
	jars := files.Match(i.ForgeDir, ".*?\\.jar")
	if len(jars) == 0 {
		return "", false
	}
	ver := getVersion(jars[0])
	return ver, ver != ""
}

func getVersion(jar string) string {
	r, e := zip.OpenReader(jar)
	if e != nil {
		return ""
	}
	defer files.Close(r)

	for _, f := range r.File {
		if f.Name == "/version.json" || f.Name == "version.json" {
			in, e := f.Open()
			if e != nil {
				continue
			}
			version := node.FromJson(in).Get("id").String()
			files.Close(in)
			return version
		}
	}

	return ""
}

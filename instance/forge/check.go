package forge

import (
	"log"
	"path/filepath"

	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/node"
)

func IsInstalled(i *modpack.Installation) bool {
	path := files.MustFile(i.GameDir, "launcher_profiles.json")
	id := node.FromJsonAt(path).Get("profiles").Get("forge").Get("lastVersionId").String()
	if id == "" {
		log.Println("lastVersionId not found")
		return false
	}

	json := filepath.Join(i.GameDir, "versions", id, id+".json")
	if !files.Exists(json) {
		log.Println("json not found")
		return false
	}

	return true
}

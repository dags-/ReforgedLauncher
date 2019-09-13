package optifine

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
)

var (
	baseUrl     = "https://optifine.net/"
	lookupUrl   = baseUrl + "adloadx?f=%s"
	hrefPattern = regexp.MustCompile(`.*?<a href=["'](downloadx.*?)["']`)
)

func Download(name, dest string, listener progress.Listener) error {
	r, e := http.Get(fmt.Sprintf(lookupUrl, name))
	if e != nil {
		return e
	}
	defer files.Close(r.Body)

	s := bufio.NewScanner(r.Body)
	for s.Scan() {
		line := s.Text()
		matches := hrefPattern.FindAllStringSubmatch(line, -1)
		for _, groups := range matches {
			if len(groups) > 1 {
				download := baseUrl + groups[1]
				return tasks.Download(download, dest, listener)
			}
		}
	}

	return fmt.Errorf("could not determine download link for: %s", name)
}

package launcher

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
)

var nameCleaner = regexp.MustCompile("[^a-zA-Z0-9 _-]")

func (l *Launcher) modpacks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		l.getModpacks(w, r)
	}
	if r.Method == http.MethodPost {
		l.postModpacks(w, r)
	}
}

func (l *Launcher) getModpacks(w http.ResponseWriter, r *http.Request) {
	modpacks, e := l.loadMopacks()
	if e == nil {
		success(w, modpacks)
	} else {
		fail(w, e)
	}
}

func (l *Launcher) postModpacks(w http.ResponseWriter, r *http.Request) {
	var repo modpack.Repo
	e := files.ParseJson(r.Body, &repo)
	if e != nil {
		fail(w, e)
		return
	}

	meta, e := repo.GetMeta()
	if e != nil {
		fail(w, e)
		return
	}

	packs, e := l.loadMopacks()
	if e != nil {
		fail(w, e)
		return
	}

	for _, m := range packs {
		if m.Repo.String() == meta.Repo.String() {
			fail(w, fmt.Errorf("modpack already added"))
			return
		}
	}

	packs = append(packs, meta)
	path := files.MustFile(l.AppDir, "modpacks.json")
	e = files.WriteJsonAt(path, packs)

	if e == nil {
		success(w, nil)
	} else {
		fail(w, e)
	}
}

func (l *Launcher) loadMopacks() ([]*modpack.Meta, error) {
	var modpacks []*modpack.Meta

	path := files.MustFile(l.AppDir, "modpacks.json")
	e := files.ParseJsonAt(path, &modpacks)
	errs.Log("Load Modpacks", e)

	if e != nil {
		modpacks, e = l.fetchModpacks()
		if e != nil {
			return nil, e
		}
		e = files.WriteJsonAt(path, &modpacks)
		errs.Log("Save Modpacks", e)
	}

	return modpacks, nil
}

func (l *Launcher) fetchModpacks() ([]*modpack.Meta, error) {
	repos := l.getRepos()

	var modpacks []*modpack.Meta
	for _, r := range repos {
		meta, e := r.GetMeta()
		if e == nil {
			modpacks = append(modpacks, meta)
		}
	}

	return modpacks, nil
}

func (l *Launcher) getRepos() []*modpack.Repo {
	var repos []*modpack.Repo
	if l.ModPacksURL != "" {
		e := tasks.GetJson(l.ModPacksURL, &repos)
		if e == nil {
			return repos
		}
	}
	return repos
}
